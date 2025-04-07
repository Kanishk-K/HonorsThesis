from BASE import DataAquisitionBase
import pandas as pd
import datetime

class CAISO(DataAquisitionBase):
    def __init__(self, start_date:datetime.datetime, end_date:datetime.datetime) -> None:
        """
        Initialize the CAISO data acquisition class.
        """
        self.dfs = []
        super().__init__(start_date, end_date)

    def _numDates(self):
        """
        Get the date range for the data.
        """
        days = (self.end_date - self.start_date).days
        return days
    
    def getData(self):
        """
        Get data from CAISO.
        """
        for i in range(self._numDates()):
            date = self.start_date + datetime.timedelta(days=i)
            url = f"https://www.caiso.com/outlook/history/{date.strftime('%Y%m%d')}/co2.csv"
            data = pd.read_csv(url)
            # Use the date provided by date and the Time column to create a DateTime column
            data['DateTime'] = pd.to_datetime(date.strftime('%Y-%m-%d') + ' ' + data['Time'])
            # Drop the Time column
            data = data.drop(columns=['Time'])
            self.dfs.append(data)
    
        self.df = pd.concat(self.dfs, ignore_index=True)
    def writeToFile(self, file_path:str) -> None:
        self.df.to_csv(file_path, index=False)
    
    def __str__(self) -> str:
        """
        Get string representation of the class.
        """
        return f"CAISO data from {self.start_date} to {self.end_date}\n {self.df}"


if __name__ == "__main__":
    start_date = datetime.datetime(2025, 1, 1)
    end_date = datetime.datetime(2025, 1, 10)
    caiso = CAISO(start_date, end_date)
    print(caiso.df.info())
    print(caiso)
