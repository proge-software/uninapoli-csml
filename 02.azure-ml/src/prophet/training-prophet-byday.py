from fbprophet import Prophet
from pandas.core.frame import DataFrame
from sklearn.model_selection import train_test_split
import pandas as pd
import numpy as np
from datetime import datetime

def parse_date(v):
    try:
        return datetime.strptime(v, "%Y-%m-%d")
    except:
        # apply whatever remedies you deem appropriate
        pass
    return v

path_to_data = "../Data/occurrencies/aggregated/occs_all.csv"

if __name__ == "__main__":    
    xf = pd.read_csv(path_to_data, date_parser=lambda x: parse_date(x), parse_dates=['Date'],
                    encoding="UTF-8", sep=',', quotechar='"', error_bad_lines=False)

    xf.rename(columns={'Date':'ds', 'Count':'y'}, inplace=True)

    print("Len xf: ", len(xf))
    xf.head()

    model = Prophet()
    model_fit = model.fit(xf)

    future = model.make_future_dataframe(periods=12, freq='D')
    future.tail()

    forecast = model.predict(future)
    forecast[['ds', 'yhat', 'yhat_lower', 'yhat_upper']].tail()

    print(forecast[:20])

    from fbprophet.diagnostics import cross_validation
    df_cv = cross_validation(model_fit, initial='3650 days', period='180 days', horizon = '365 days')
    df_cv.to_csv("crossvalid.csv")
    print(df_cv)

    from fbprophet.diagnostics import performance_metrics
    df_p = performance_metrics(df_cv)
    df_p.head()
    df_cv.to_csv("out-pmetrics.csv")
