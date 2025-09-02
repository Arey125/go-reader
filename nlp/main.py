from fastapi import FastAPI
from pydantic import BaseModel
import spacy
from spacy.tokens import Token
 
app = FastAPI()

nlp = spacy.load("en_core_web_sm")


class Text(BaseModel):
    content: str


def token_filter(token: Token):
    if token.pos_ not in ['ADJ', 'ADV', 'INTJ', 'NOUN', 'VERB']:
        return False
    if token.is_stop:
        return False
    if not token.is_alpha:
        return False
    return True


@app.post("/")
async def root(text: Text):
    doc = nlp(text.content)
    result = [
        {
            "text": token.text,
            "lemma": token.lemma_,
            "pos": token.pos_,
            "start": token.idx,
            "end": token.idx + len(token.text),
        }
        for token in doc
        if token_filter(token)
    ]
    return {"result": result}
