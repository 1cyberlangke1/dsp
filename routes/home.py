# -*- coding: utf-8 -*-
"""首页路由"""
from fastapi import APIRouter, Request
from fastapi.templating import Jinja2Templates

from core.config import TEMPLATES_DIR

router = APIRouter()
templates = Jinja2Templates(directory=TEMPLATES_DIR)


@router.get("/")
def index(request: Request):
    return templates.TemplateResponse("welcome.html", {"request": request})
