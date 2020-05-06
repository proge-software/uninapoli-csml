FROM python

WORKDIR /workspace
COPY occurrencies_byday.py src/

RUN mkdir -p /workspace/out \
    && pip install --upgrade tqdm pandas azureml-sdk

WORKDIR /workspace/src
CMD ["python", "occurrencies_byday.py"]
