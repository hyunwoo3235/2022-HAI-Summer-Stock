import argparse
import zipfile

import h5py
import pandas as pd

from utils import logging

logger = logging.get_logger(__name__)


def parse_time(time: str):
    return pd.to_datetime(f"{time}", format="%H%M%S")


def time_interval(prev_time, time):
    return (parse_time(time) - parse_time(prev_time)).total_seconds()


def convert_stooq_to_dataset(stooq_file: str, dataset_file: str):
    hp = h5py.File(dataset_file, "w")
    zf = zipfile.ZipFile(stooq_file, "r")
    files = [f for f in zf.namelist() if f.endswith(".txt")]
    for file in tqdm(files):
        try:
            df = pd.read_csv(zf.open(file))
        except Exception as e:
            logger.error(f"Error reading file {file}: {e}")
            continue

        prev_date = df.iloc[0]["<DATE>"]
        prev_time = df.iloc[0]["<TIME>"]
        new_df = pd.DataFrame(columns=["date", "time", "price"])
        for i in range(len(df)):
            date = df.iloc[i]["<DATE>"]
            time = df.iloc[i]["<TIME>"]
            price = df.iloc[i]["<OPEN>"] * 100

            if prev_date == date and time_interval(prev_time, time) != 600:
                continue

            prev_date = date
            prev_time = time

            new_df = new_df.append({"date": date, "time": time, "price": price}, ignore_index=True)

        hp[file] = new_df.astype(int).values


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--stooq_file", type=str, help="Path to the zipfile downloaded from Stooq")
    parser.add_argument("--dataset_file", default="dataset.hdf5", type=str, help="Path to the dataset file")
    args = parser.parse_args()
    convert_stooq_to_dataset(args.stooq_file, args.dataset_file)
