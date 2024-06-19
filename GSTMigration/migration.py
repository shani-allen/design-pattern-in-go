import csv
import requests
from requests.adapters import HTTPAdapter
from requests.packages.urllib3.util.retry import Retry

# Define constants
CSV_FILE_PATH = 'receipts.csv'  # Replace with your CSV file path
URL = 'http://localhost:8002/v1/ecomm-workflow/order/document'  # Replace with the actual URL
DOC_TYPE = 'INVOICE'
TENANT_ID = 'aUSsW8GI03dyQ0AIFVn92'

# Function to read receipt IDs from CSV
def read_receipt_ids(csv_file_path):
    receipt_ids = []
    with open(csv_file_path, mode='r') as file:
        csv_reader = csv.reader(file)
        # header = next(csv_reader)  # Skip the header if there is one
        for row in csv_reader:
            if row:  # Ensure row is not empty
                receipt_ids.append(row[0])
    return receipt_ids

# Function to make the HTTP request
def make_request(session, receipt_id):
    payload = {
        "doc_type": DOC_TYPE,
        "receipt_id": receipt_id,
        "tenant_id": TENANT_ID
    }
    try:
        response = session.post(URL, json=payload)
        response.raise_for_status()  # Raise an HTTPError for bad responses
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
    receipt_ids = read_receipt_ids(CSV_FILE_PATH)
    print("receipt_ids:", receipt_ids)

    # Set up retry strategy
    retry_strategy = Retry(
        total=3,  # Total number of retries
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
            print(response.json())  # Print response JSON, adjust as needed

if __name__ == "__main__":
    main()
