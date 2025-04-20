from abc import ABC, abstractmethod
import pandas as pd
from matplotlib import pyplot as plt
import os

class AbstractIntensityViz(ABC):
    def __init__(self, filepath:str) -> None:
        pass

    @abstractmethod
    def plot(self, ax=None) -> None:
        pass

class BaseIntensityViz(AbstractIntensityViz):
    def __init__(self, filepath:str) -> None:
        """
        Initialize the base intensity visualization class with a file path.
        """
        self.filepath = filepath
        self.df = pd.read_csv(filepath)
        if self.df.empty:
            raise ValueError("Dataframe is empty.")
        if not isinstance(self.df, pd.DataFrame):
            raise ValueError("Dataframe is not a pandas dataframe.")
    
    def plot(self, ax=None) -> None:
        """
        Plot the data.
        """
        # Convert start_date to datetime
        self.df['start_date'] = pd.to_datetime(self.df['start_date'])
        # Group by hour and minute
        self.df['hour_minute'] = self.df['start_date'].dt.strftime('%H:%M')
        # Average "generated_rate_kg_per_mwh" by this group
        self.aggregated_df = self.df.groupby('hour_minute').mean().reset_index()
        # Plot the data
        self.aggregated_df.plot(x='hour_minute', y='generated_rate_kg_per_mwh', kind='line', ax=ax)
        plt.xlabel('Hour of the Day')
        plt.ylabel('Average Generated Rate (kg/MWh)')
        plt.title(f"Average Generated Rate by Hour of the Day ({self.filepath.replace('collected/', '')[:-4]})")

class AllIntensityViz(AbstractIntensityViz):
    def __init__(self, filepath:str) -> None:
        # Treat the filepath as a directory
        self.filepath = filepath
    
    def plot(self):
        # Get all CSV files in the directory, create a list of BaseIntensityViz objects
        files = [f for f in os.listdir(self.filepath) if f.endswith('.csv')]
        intensity_viz_objects = [BaseIntensityViz(os.path.join(self.filepath, file)) for file in files]
        # Plot each file on the same graph
        fig, ax = plt.subplots()
        for intensity_viz in intensity_viz_objects:
            intensity_viz.plot(ax=ax)

        ax.set_title("Average Generated Rate by Hour of the Day (All ISOs)")
        # Add legend
        labels = [os.path.basename(intensity_viz.filepath).replace('.csv', '') for intensity_viz in intensity_viz_objects]
        ax.legend(labels, loc='upper left', bbox_to_anchor=(1, 1))


if __name__ == "__main__":
    # Example usage (single)

    # filepath = "collected/MISO.csv"
    # intensity_viz = BaseIntensityViz(filepath)
    # intensity_viz.plot()

    # Example usage (all)
    filepath = "collected/"
    intensity_viz = AllIntensityViz(filepath)
    intensity_viz.plot()

    plt.show()