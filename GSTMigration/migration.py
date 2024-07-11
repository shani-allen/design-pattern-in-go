import csv
import requests
from requests.adapters import HTTPAdapter
from requests.packages.urllib3.util.retry import Retry
import logging
import yaml
import json
from typing import Any
from Crypto.Cipher import AES
from Crypto.Util.Padding import pad
import base64
from urllib.parse import urlencode

# Define constants
URL = 'http://localhost:8002/v1/ecomm-workflow/order/document'
DOC_TYPE = 'INVOICE'
TENANT_ID = 'aUSsW8GI03dyQ0AIFVn92'
CONFIG_PATH = 'config.yaml'
CSV_FILE_PATH = 'form.csv'

# AES encryption function
def get_aes_encrypted(plaintext: str, secret_key: str, secret_iv: str) -> str:
    try:
        # Convert plaintext, key, and IV to bytes
        plaintext_bytes = plaintext.encode('utf-8')
        key_bytes = secret_key.encode('utf-8')
        iv_bytes = secret_iv.encode('utf-8')

        # Pad plaintext to be a multiple of the block size
        block_size = 16
        padded_plaintext = pad(plaintext_bytes, block_size)

        # Create AES cipher
        cipher = AES.new(key_bytes, AES.MODE_CBC, iv_bytes)

        # Encrypt the padded plaintext
        ciphertext = cipher.encrypt(padded_plaintext)

        # Encode the ciphertext in base64
        ciphertext_base64 = base64.b64encode(ciphertext).decode('utf-8')

        return ciphertext_base64

    except Exception as e:
        error_msg = f"AES Encryption Error, val: {plaintext}, err: {str(e)}"
        raise Exception(error_msg)

# Load the configuration from the YAML file
def load_config(config_path: str) -> dict:
    with open(config_path, 'r') as f:
        config_data = yaml.safe_load(f)
    return config_data

class AllenErpAPIImpl:
    def __init__(self, config: dict, log: logging.Logger):
        self.config = config
        self.log = log
        self.client = requests.Session()

    def create_access_token(self, ctx: Any, token_request: Any) -> dict:
        user_name = self.config['username']
        password = self.config['password']
        grant_type = "password"

        key = f"{self.config['client_id']}:{self.config['client_secret']}"
        encrypted_key = base64.b64encode(key.encode()).decode()

        data = {
            "username": user_name,
            "password": password,
            "grant_type": grant_type
        }

        url_encoded_data = urlencode(data)
        endpoint = f"{self.config['url']}{self.config['middleware_route']}/token"

        headers = {
            "Authorization": f"BASIC {encrypted_key}",
            "Content-Type": "application/x-www-form-urlencoded"
        }

        try:
            response = self.client.post(endpoint, data=url_encoded_data, headers=headers)
            response.raise_for_status()
        except requests.RequestException as e:
            self.log.error(f"ALLEN ERP: error while creating request for CreateSession API, url: {endpoint}, err: {str(e)}")
            raise Exception(f"Error while creating request for CreateSession API, err: {str(e)}")

        try:
            token_response = response.json()
        except json.JSONDecodeError as e:
            self.log.error(str(e))
            raise

        if not token_response.get('access_token'):
            error_msg = "Error in fetching access token from ERP"
            self.log.error(error_msg)
            raise Exception(error_msg)

        return token_response

    def get_receipt_details_by_receipt_id(self, ctx: Any, access_token: str, request: Any) -> dict:
        try:
            json_data = json.dumps(request)
        except TypeError as e:
            self.log.error(str(e))
            raise

        encrypted_data = get_aes_encrypted(json_data, self.config['aes_secret_key'], self.config['aes_secret_iv'])
        root = get_aes_encrypted(self.config['get_receipt_details_by_receipt_id']['get_receipt']['root'],
                                 self.config['aes_secret_key'], self.config['aes_secret_iv'])
        app = get_aes_encrypted(self.config['get_receipt_details_by_receipt_id']['get_receipt']['app'],
                                self.config['aes_secret_key'], self.config['aes_secret_iv'])
        session = get_aes_encrypted("2024-2025",
                                    self.config['aes_secret_key'], self.config['aes_secret_iv'])
        data = {"encryptdata": encrypted_data}

        url_encoded_data = urlencode(data)

        endpoint = f"{self.config['url']}{self.config['middleware_route']}{self.config['get_receipt_details_by_receipt_id']['get_receipt']['api_route']}"

        headers = {
            "Authorization": f"Bearer {access_token}",
            "Root": root,
            "App": app,
            "Session": session,
            "Content-Type": "application/x-www-form-urlencoded"
        }

        try:
            response = self.client.post(url=endpoint, data=url_encoded_data, headers=headers)
            print("\n",response.content)
            response.raise_for_status()
        except requests.RequestException as e:
            self.log.error(f"ALLEN ERP: error while creating request for GetReceiptDetailsByReceiptID API, url: {endpoint}, err: {str(e)}")
            raise Exception(f"Error while creating request for GetReceiptDetailsByReceiptID API, err: {str(e)}")

        return self.format_get_erp_response(ctx, response.content)

    def format_get_erp_response(self, ctx: Any, resp_body: bytes) -> dict:
        try:
            init_resp = json.loads(resp_body)
        except json.JSONDecodeError as e:
            self.log.error(str(e))
            raise
        print("message:",init_resp["message"])
        if isinstance(init_resp['result'], str) and init_resp['status'] == "FAILURE":
            error_data = init_resp['result']
            error_response = {
                'message': init_resp['message'],
                'status': init_resp['status'],
                'error': {'message': error_data}
            }
            return error_response

        print(init_resp['result'])
        if not isinstance(init_resp['result'], dict):
            self.log.error("Failed to convert Result to a map")
            # raise Exception("Failed to convert Result to a map")

        erp_get_receipt_response = {
            'message': init_resp['message'],
            'status': init_resp['status'],
            'result': {'GSTFeeReceipt': init_resp['result'].get('GST Fee Receipt', [])}
        }

        return erp_get_receipt_response

# Function to make the HTTP request
def make_request(session, receipt_id):
    payload = {
        "doc_type": DOC_TYPE,
        "receipt_id": receipt_id,
        "tenant_id": TENANT_ID
    }
    print("ecomm-payload:",payload)
    try:
        response = session.post(URL, json=payload)
        response.raise_for_status()
        return response
    except requests.exceptions.HTTPError as errh:
        print(f"HTTP error occurred: {errh}")
    except requests.exceptions.ConnectionError as errc:
        print(f"Connection error occurred: {errc}")
    except requests.exceptions.Timeout as errt:
        print(f"Timeout error occurred: {errt}")
    except requests.exceptions.RequestException as err:
        print(f"An error occurred: {err}")

# Main script execution
def main():
    # Load configuration
    config = load_config(CONFIG_PATH)

    # Set up logging
    logging.basicConfig(level=logging.INFO)
    logger = logging.getLogger(__name__)

    # Initialize the API implementation
    api_impl = AllenErpAPIImpl(config, logger)

    # Retrieve access token
    req={
         "username": "",
        "password": "",
        "grant_type": ""
    }

    access_token_response = api_impl.create_access_token(None, req)
    access_token = access_token_response.get('access_token')
    print(access_token)
    # Read form IDs from CSV
    form_ids = read_form_ids_from_csv(CSV_FILE_PATH)

    # Process each form ID to get receipt details
    for form_id in form_ids:
        request_data = {
            "fno": form_id
        }
        print(form_id)
        receipt_details_response = api_impl.get_receipt_details_by_receipt_id(None, access_token, request_data)
        receipt_ids = [receipt['receiptnumber'] for receipt in receipt_details_response['result']['GSTFeeReceipt']]

        # Set up retry strategy for eComm workflow request
        retry_strategy = Retry(
            total=1,  # Total number of retries
            backoff_factor=1,  # Wait between retries (exponential backoff)
            status_forcelist=[429, 500, 502, 503, 504],  # HTTP status codes to retry on
            allowed_methods=["POST"],  # Methods to retry
        )

        adapter = HTTPAdapter(max_retries=retry_strategy)
        session = requests.Session()
        session.mount("http://", adapter)
        session.mount("https://", adapter)

        for receipt_id in receipt_ids:
            response = make_request(session, receipt_id)
            if response:
                print(f"Request for Receipt ID {receipt_id} returned status code: {response.status_code}")
                print(response.json())

def read_form_ids_from_csv(csv_file_path):
    form_ids = []
    with open(csv_file_path, mode='r') as file:
        csv_reader = csv.reader(file)
        for row in csv_reader:
            if row:  # Ensure row is not empty
                form_ids.append(row[0])
    return form_ids

if __name__ == "__main__":
    main()
