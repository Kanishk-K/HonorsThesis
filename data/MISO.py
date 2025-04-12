from BASE import DataAquisitionBase
import datetime
from dotenv import load_dotenv

class DataAquisitionMISO(DataAquisitionBase):
    """
    Implementation of the data acquisition class for the Midcontinent Independent System Operator.
    """
    def __init__(self, start_date:datetime.datetime, end_date:datetime.datetime, ISO:str) -> None:
        super().__init__(start_date, end_date, ISO)

if __name__ == "__main__":
    load_dotenv()
    # Example usage
    start_date = datetime.datetime(2024, 1, 1)
    end_date = datetime.datetime(2025, 1, 1)
    
    data_acquisition = DataAquisitionMISO(start_date, end_date, "MISO")
    data_acquisition.writeToFile("collected/MISO.csv")