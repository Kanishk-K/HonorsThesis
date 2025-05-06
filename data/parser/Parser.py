from abc import ABC, abstractmethod
from typing import List, Dict
import pandas as pd
import re

class AbstractParser(ABC):
    def __init__(self, filepath:str) -> None:
        """
        Initialize the base parser class with a file path.
        """
        self.filepath = filepath
        self.df = None
        self._parse()

    @abstractmethod
    def _parse(self) -> None:
        """
        Abstract method to parse the data.
        """
        pass

    def get_df(self) -> pd.DataFrame:
        """
        Get the parsed dataframe.
        """
        if self.df is None:
            raise ValueError("Dataframe is not parsed yet.")
        return self.df
    
class BaseParser(AbstractParser):
    def _parseMapCarbon(map_str:str) -> Dict[str, float]:
        entries = re.findall(r'\{(\w+)[^}]*\}:(-?\d+(?:\.\d+)?(?:[eE][+-]?\d+)?)+', map_str)
        return {k: float(v) for k, v in entries}

    def _parseMapSLO(map_str:str) -> Dict[str, int]:
        entries = re.findall(r'\{(\w+)[^}]*\}:(\d+)', map_str)
        return {k: int(v) for k, v in entries}

    def _parse(self) -> None:
        """
        Parse the data from the file.
        """
        with open(self.filepath, 'r') as file:
            text = file.read()
        
        entries = re.findall(r'Simulator State:.*?Completed Jobs Length: \d+', text, re.DOTALL)
        dfVals = []
        # Carbon Emissions
        print(len(entries))
        for entry in entries:
            carbonMatch = re.search(r"Carbon Emission:\s+map\[(.*?)\]", entry)
            sloMatch = re.search(r"SLO Timeouts:\s+map\[(.*?)\]", entry)
            carbonEntries = BaseParser._parseMapCarbon(carbonMatch.group(1)) if carbonMatch else {}
            sloEntries = BaseParser._parseMapSLO(sloMatch.group(1)) if sloMatch else {}
            dataVals = []
            for key, value in carbonEntries.items():
                dataVals.append((key, value, sloEntries.get(key, 0)))
            
            dfVals.extend(dataVals)

        # Create DataFrame
        self.df = pd.DataFrame(dfVals, columns=['Job', 'Carbon Emission', 'SLO Timeout'])

if __name__ == "__main__":
    # Example usage
    parser = BaseParser('summary_MISO_6hr_random_hybridSelection_80.log')
    df = parser.get_df()
    print(df["Carbon Emission"].mean())
    print(df["SLO Timeout"].mean())