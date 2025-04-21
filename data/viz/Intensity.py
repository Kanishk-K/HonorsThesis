from abc import ABC, abstractmethod
import pandas as pd
from matplotlib import pyplot as plt
from matplotlib.dates import DateFormatter
from matplotlib.collections import PolyCollection
import seaborn as sns
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
    
    def plot(self, ax=None, label=None) -> None:
        """
        Plot the data.
        """
        self.df['start_date'] = pd.to_datetime(self.df['start_date'])
        self.df['hour_minute'] = pd.to_datetime(self.df['start_date'].dt.strftime('%H:%M'), format='%H:%M')

        if ax is None:
            fig, ax = plt.subplots()
            # Set the title of the plot
            ISOName = os.path.basename(self.filepath).replace('.csv', '')
            ax.set_title(f"Average Generated Rate by Hour of the Day ({ISOName}) w/ SD")


        # Plot the line with label
        sns.lineplot(
            data=self.df,
            x='hour_minute',
            y='generated_rate_kg_per_mwh',
            ax=ax,
            errorbar='sd',
            label=label  # pass label from parent
        )

        ax.xaxis.set_major_formatter(DateFormatter("%H:%M"))
        ax.set_xlabel('Hour of the Day')
        ax.set_ylabel('Average Generated Rate (kg CO₂/MWh)')

class AllIntensityViz(AbstractIntensityViz):
    def __init__(self, filepath:str) -> None:
        # Treat the filepath as a directory
        self.filepath = filepath
    
    def plot(self):
        files = [f for f in os.listdir(self.filepath) if f.endswith('.csv')]
        intensity_viz_objects = [BaseIntensityViz(os.path.join(self.filepath, file)) for file in files]

        fig, ax = plt.subplots()
        for intensity_viz in intensity_viz_objects:
            label = os.path.basename(intensity_viz.filepath).replace('.csv', '')
            intensity_viz.plot(ax=ax, label=label)

        ax.set_title("Average Generated Rate by Hour of the Day (All ISOs)")

        # Fix legend to remove CI bands
        handles, labels = ax.get_legend_handles_labels()
        filtered = [(h, l) for h, l in zip(handles, labels) if not isinstance(h, PolyCollection)]
        if filtered:
            ax.legend(*zip(*filtered), loc='upper left', bbox_to_anchor=(1, 1))


if __name__ == "__main__":
    # Example usage (single)

    # filepath = "collected/CAISO.csv"
    # intensity_viz = BaseIntensityViz(filepath)
    # intensity_viz.plot()

    # Example usage (all)
    filepath = "collected/"
    intensity_viz = AllIntensityViz(filepath)
    intensity_viz.plot()

    plt.show()