FROM python:3-alpine

ENV FLASK_APP=server.py

WORKDIR /code

COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt

COPY . .

CMD [ "python", "-m", "flask", "run", "--host=0.0.0.0", "--port=80" ]
