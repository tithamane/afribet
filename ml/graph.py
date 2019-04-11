"""Generate graphs based on files saved in the data folder."""

import os
import pandas as pd
import matplotlib.pyplot as plt

from datetime import datetime

class SevenBallColors():
    yellow_numbers = [1, 3, 5, 8, 10, 12, 13, 15, 17, 20, 22, 24, 26, 27, 29, 32, 34, 36, 37, 39, 41]

    @staticmethod
    def is_yellow(number: int):
        yellow = number in SevenBallColors.yellow_numbers
        return 1 if yellow else 0

class SevenBallTotal():
    def __init__(self, df):
        self.df = df


def more_yellows(row):
    return 1 if row['yellows'] > 3 else 0

if __name__ == '__main__':
    print("It works")
    file_path = os.path.join(os.getcwd(), '..', 'scraper', 'data', 'csv', 'seven_ball.csv')
    df = None
    df = pd.read_csv(file_path)
    
    print(df[0:5])

    # Create date from unix timestamp
    df['date'] = df['unix_time'].apply(lambda date: datetime.utcfromtimestamp(date).strftime("%Y-%m-%d"))

    df_all_total = df.copy()
    df_all_colors = df.copy()

    columns = [
        'first_number',
        'second_number',
        'third_number',
        'fourth_number',
        'fifth_number',
        'sixth_number',
        'seventh_number'
    ]

    # Create a column with a total sum of column values
    sum_columns = 0
    color_columns = 0
    for column_name in columns:
        # Add column
        sum_columns += df[column_name]

        # Set colors for each column
        df_all_colors[column_name] =  df_all_colors[column_name].apply(SevenBallColors.is_yellow)

        # Count yellow colors
        color_columns += df_all_colors[column_name]
         

    # Compute total
    df_all_total['all_total'] = sum_columns
    df_totals = df_all_total[['date', 'all_total']]

    # Compute colors
    df_all_colors['yellows'] = color_columns
    df_all_colors['blacks'] = 7 - df_all_colors['yellows']

    df_all_colors['more_yellows'] = df_all_colors.apply(more_yellows, axis=1)

    df_colors = df_all_colors[['date', 'yellows', 'blacks', 'more_yellows']]

    print(df_totals[0:5])
    print('---------------')
    print(df_colors[0:5])

