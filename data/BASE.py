from abc import ABC, abstractmethod
import datetime
import os
import requests
import pandas as pd

class AbstractDataAquisition(ABC):
    """
    Describes the abstract class for data acquisition.
    """

    def __init__(self, start_date:datetime.datetime, end_date:datetime.datetime, ISO:str) -> None:
        pass

    @abstractmethod
    def getData(self) -> None:
        """
        Calls the API to get the data.
        """
        pass
    
    @abstractmethod
    def writeToFile(self, file_path:str) -> None:
        """
        Writes the data to a file.
        """
        pass

    @abstractmethod
    def __str__(self) -> str:
        """
        Returns the string representation of the class.
        """
        pass


class DataAquisitionBase(AbstractDataAquisition):
    """
    Basic implementation of the data acquisition class.
    """
    def __init__(self, start_date:datetime.datetime, end_date:datetime.datetime, ISO:str) -> None:
        """
        Initialize the data acquisition class with a file path.
        """
        self.start_date = start_date
        self.end_date = end_date
        self.ISO = ISO
        if os.getenv("SINGULARITY_API") is None:
            raise ValueError("SINGULARITY_API environment variable not set.")
        if os.getenv("SINGULARITY_API") == "":
            raise ValueError("SINGULARITY_API environment variable is empty.")
        if start_date > end_date:
            raise ValueError("Start date must be less than end date.")
        if end_date > datetime.datetime.now():
            raise ValueError("End date must be less than current date.")
        self.getData()

    def _numDates(self):
        """
        Get the date range for the data.
        """
        days = (self.end_date - self.start_date).days
        return days

    def getData(self) -> None:
        dfs = []
        for i in range(self._numDates()):
            date = self.start_date + datetime.timedelta(days=i)
            print(f"Getting data for {date.strftime('%Y-%m-%d')}")
            params = {
                "emission_factor_adjustment": "for_electricity_adjusted",
                "pollutant": "co2",
                "resolution": "5m",
                # Datetimes in ISO-8601 format (2022-01-07T00:00:00Z)
                "start": date.isoformat() + "Z",
                "end": (date + datetime.timedelta(days=1)).isoformat() + "Z",
                "region": self.ISO,
                "source": "ISO",
                "emission_factor_source": "EGRID"
            }
            headers = {
                "x-api-key": os.getenv("SINGULARITY_API"),
            }
            response = requests.get(
                "https://api.singularity.energy/v2/generated/carbon-intensity",
                params=params,
                headers=headers
            )
            if response.status_code != 200:
                raise ValueError(f"Error: {response.status_code} - {response.text}")
            data = response.json()
            if data["data"] is None:
                raise ValueError(f"Error: {response.status_code} - {response.text}")
            date_df = [(
                d["start_date"],
                d["data"]["generated_rate_kg_per_mwh"]
            ) for d in data["data"]]
            dfs.append(pd.DataFrame(date_df, columns=["start_date", "generated_rate_kg_per_mwh"]))
        self.df = pd.concat(dfs, ignore_index=True)
        self.df["ISO"] = self.ISO

    
    def writeToFile(self, file_path):
        self.df.to_csv(file_path, index=False)

    def __str__(self) -> str:
        """
        Returns the string representation of the class.
        """
        return f"Data for {self.ISO} from {self.start_date} to {self.end_date}\n{self.df}"