FROM python:3.10-alpine3.17

ADD src /app

WORKDIR /app

RUN pip install -i https://pypi.tuna.tsinghua.edu.cn/simple -r requirements.txt

RUN apk add nodejs npm

ENTRYPOINT ["python3","PackerFuzzer.py","-s","true"]