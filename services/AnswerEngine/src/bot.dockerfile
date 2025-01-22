FROM python:3.11
LABEL authors="kirillgorunov"

COPY ./ /app/AnswerEngine

WORKDIR /app

RUN pip install -r AnswerEngine/src/requirements.txt

CMD ["python3", "AnswerEngine/src/main.py"]

ENV PYTHONPATH=/app

