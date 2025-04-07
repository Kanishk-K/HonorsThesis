from abc import ABC, abstractmethod
import datetime

class DataAquisitionBase(ABC):
    """
    Abstract class for data acquisition.
    """
    def __init__(self, start_date:datetime.datetime, end_date:datetime.datetime) -> None:
        """
        Initialize the data acquisition class with a file path.
        """
        self.start_date = start_date
        self.end_date = end_date
        if start_date > end_date:
            raise ValueError("Start date must be less than end date.")
        if end_date > datetime.datetime.now():
            raise ValueError("End date must be less than current date.")
        self.getData()

    @abstractmethod
    def getData(self) -> None:
        """
        Abstract method to get data.
        """
        pass

    @abstractmethod
    def writeToFile(self, file_path:str) -> None:
        """
        Abstract method to get data frame.
        """
        pass

    @abstractmethod
    def __str__(self) -> str:
        """
        Abstract method to get string representation of the class.
        """
        pass