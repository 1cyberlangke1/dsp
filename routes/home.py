# -*- coding: utf-8 -*-
"""首页和 WebUI 路由"""
import os
from fastapi import APIRouter, Request
from fastapi.responses import HTMLResponse, FileResponse
from fastapi.templating import Jinja2Templates

from core.config import TEMPLATES_DIR, STATIC_ADMIN_DIR

router = APIRouter()
templates = Jinja2Templates(directory=TEMPLATES_DIR)


@router.get("/")
def index(request: Request):
    return templates.TemplateResponse("welcome.html", {"request": request})


@router.get("/webui")
@router.get("/webui/{path:path}")
async def webui(request: Request, path: str = ""):
    """提供 WebUI 静态文件"""
    # 检查 static/admin 目录是否存在
    if not os.path.isdir(STATIC_ADMIN_DIR):
        return HTMLResponse(
            content="<h1>WebUI not built</h1><p>Run <code>cd webui && npm run build</code> first.</p>",
            status_code=404
        )
    
    # 如果请求的是具体文件（如 js, css）
    if path and "." in path:
        file_path = os.path.join(STATIC_ADMIN_DIR, path)
        if os.path.isfile(file_path):
            return FileResponse(file_path)
        return HTMLResponse(content="Not Found", status_code=404)
    
    # 否则返回 index.html（SPA 路由）
    index_path = os.path.join(STATIC_ADMIN_DIR, "index.html")
    if os.path.isfile(index_path):
        return FileResponse(index_path)
    
    return HTMLResponse(content="index.html not found", status_code=404)
