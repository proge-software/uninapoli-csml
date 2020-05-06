from azureml.core.authentication import ServicePrincipalAuthentication
from azureml.train.automl import AutoMLConfig
from azureml.core.experiment import Experiment
from azureml.core.workspace import Workspace
from azureml.core.datastore import Datastore
import os
import numpy as np
import logging
import azureml.core
from azureml.core.compute import AmlCompute
from azureml.core.compute import ComputeTarget
from azureml.core.dataset import Dataset
from tqdm import tqdm
import pandas as pd
import os.path
from datetime import datetime, timedelta
from typing import Tuple, Dict, List
from multiprocessing import Pool
import pprint as pp
import csv

cols = ["GroCode", "fTransmis", "FuelType", "DateOut"]

# defaults
default_location="westeurope"
default_sku="basic"
default_blob_datastore_name = 'ds_customerchoice'  # Name of the datastore to workspace
default_tenant_id = ""
default_service_principal_id= ""
default_subscription_id = ""
default_resource_group = ""
default_workspace_name = ""
default_auth_secret = ""
default_container_name = ""
default_account_name = ""
default_account_key = ""
default_sas_token = ""
default_num_threads = 4

# data
cols = ['DateOut', 'DateIn', 'GroCode', 'FuelType', 'fTransmis']
dateIntervals = [["2000-01-01", "2001-01-01"],
                 ["2001-01-01", "2002-01-01"],
                 ["2002-01-01", "2003-01-01"],
                 ["2003-01-01", "2004-01-01"],
                 ["2004-01-01", "2005-01-01"],
                 ["2005-01-01", "2006-01-01"],
                 ["2006-01-01", "2007-01-01"],
                 ["2007-01-01", "2008-01-01"],
                 ["2008-01-01", "2009-01-01"],
                 ["2009-01-01", "2010-01-01"],
                 ["2010-01-01", "2011-01-01"],
                 ["2011-01-01", "2012-01-01"],
                 ["2012-01-01", "2013-01-01"],
                 ["2013-01-01", "2014-01-01"],
                 ["2014-01-01", "2015-01-01"],
                 ["2015-01-01", "2016-01-01"]]

subscription_id = os.getenv("SUBSCRIPTION_ID", default_subscription_id)
resource_group = os.getenv("RESOURCE_GROUP", default_resource_group)
workspace_name = os.getenv("WORKSPACE_NAME", default_workspace_name)
auth_secret = os.getenv("AUTH_SECRET", default_auth_secret)
blob_datastore_name = os.getenv("BLOB_DATASTORE_NAME", default_blob_datastore_name)
container_name = os.getenv("BLOB_CONTAINER", default_container_name)
account_name = os.getenv("BLOB_ACCOUNTNAME", default_account_name)
account_key = os.getenv("BLOB_ACCOUNT_KEY", default_account_key)  # Storage account key
sas_token = os.getenv("SAS_TOKEN", default_sas_token)
tenant_id = os.getenv("SAS_TOKEN", default_tenant_id)
service_principal_id = os.getenv("SERVICE_PRINCIPAL_ID", default_service_principal_id)
location = os.getenv("LOCATION", default_location)
sku = os.getenv("SKU", default_sku)
num_thread = os.getenv("THREADS", default_num_threads)


def _byday(dates: List[str]) -> Dict[str, int]:
    print("Run for dates: ", dates)
    mindate, maxdate = datetime.strptime(
        dates[0], "%Y-%m-%d"), datetime.strptime(dates[1], "%Y-%m-%d")

    gf = df.drop(df[df.DateOut > maxdate].index, inplace=False)
    gf.drop(gf[gf.DateIn < mindate].index, inplace=True)
    gf.dropna()

    print("Filtered data:", len(gf))

    occs = {}
    for _, row in gf.iterrows():
        if (row.DateOut >= mindate or row.DateOut < maxdate) and (row.DateOut < mindate or row.DateIn > mindate):
            tdelta = row.DateIn - row.DateOut  # type: timedelta

            if tdelta.days > 0:
                for i in range(0, int(tdelta.days), 1):
                    dd = row.DateOut + timedelta(days=i)
                    if dd < maxdate and dd > mindate:
                        k = dd.strftime("%Y-%m-%d")
                        occs[k] = occs.setdefault(k, 0) + 1

    dict_to_csv("occs-"+mindate.strftime("%Y"), occs)

def occs_byday(mt=6):
    p = Pool(mt)
    p.map(_byday, dateIntervals)

def dict_to_csv(fname: str, _occs: Dict[datetime, int]):
    with open(fname+".csv", mode="w") as f:
        w = csv.writer(f, delimiter=',', quotechar='"',
                       quoting=csv.QUOTE_MINIMAL)

        w.writerow(["Date", "Count"])
        for k, v in tqdm(_occs.items(), desc=fname):
            w.writerow([k, v])

def parse_date(v):
    try:
        return datetime.strptime(v, "%Y-%m-%d %H:%m:%S")
    except:
        # apply whatever remedies you deem appropriate
        pass
    return v


if __name__ == "__main__":
    # Retrieve data from Azure
    spa = ServicePrincipalAuthentication(tenant_id=tenant_id,
                                        service_principal_id=service_principal_id, 
                                        service_principal_password=auth_secret)

    ws = Workspace(subscription_id, resource_group, workspace_name, auth=spa,
                _location="westeurope", _disable_service_check=False, _workspace_id=None, sku='basic')

    datastore = ws.get_default_datastore()
    datastore.download("../Data/", prefix="CSV/")

    # Read downloaded data
    df = pd.read_csv("../Data/CSV/extraction.csv", date_parser=lambda x: parse_date(x), usecols=cols, parse_dates=[
                     'DateOut', 'DateIn'], encoding="UTF-16 LE", sep=';', quotechar='"', error_bad_lines=False)
    print("Read data", len(df))

    # Filter data
    df["fTransmis"] = df["fTransmis"].transform(
        lambda x: "M" if x == "S" else "A")
    print("Transformed transmis data", len(df))

    # TODO: Drop where fueltype is E or H
    df.drop(df[df.FuelType == "E"].index, inplace=True)
    print("dropped fueltype E", len(df))

    df.drop(df[df.FuelType == "H"].index, inplace=True)
    print("dropped fueltype H", len(df))

    # Data before 2000 are not valuable, removing cars returned before 2000
    # i.e. in 2000 were occupied
    df.drop(df[df.DateIn < datetime(2000, 1, 1)].index, inplace=True)
    print("dropped datein before 2000-01-01", len(df))

    df = df.dropna()
    print("dropped na; filtered data", len(df))

    # Process data
    occs_byday(mt=num_thread)
