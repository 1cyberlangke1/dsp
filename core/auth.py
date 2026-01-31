# -*- coding: utf-8 -*-
"""账号认证与管理模块"""
import random
from fastapi import HTTPException, Request

from .config import CONFIG, logger
from .deepseek import login_deepseek_via_account, BASE_HEADERS

# -------------------------- 全局账号队列 --------------------------
account_queue = []  # 维护所有可用账号
claude_api_key_queue = []  # 维护所有可用的Claude API keys


def init_account_queue():
    """初始化时从配置加载账号"""
    global account_queue
    account_queue = CONFIG.get("accounts", [])[:]  # 深拷贝
    random.shuffle(account_queue)  # 初始随机排序


def init_claude_api_key_queue():
    """Claude API keys由用户自己的token提供，这里初始化为空"""
    global claude_api_key_queue
    claude_api_key_queue = []


# 初始化
init_account_queue()
init_claude_api_key_queue()


# ----------------------------------------------------------------------
# 辅助函数：获取账号唯一标识（优先 email，否则 mobile）
# ----------------------------------------------------------------------
def get_account_identifier(account: dict) -> str:
    """返回账号的唯一标识，优先使用 email，否则使用 mobile"""
    return account.get("email", "").strip() or account.get("mobile", "").strip()


# ----------------------------------------------------------------------
# 账号选择与释放
# ----------------------------------------------------------------------
def choose_new_account(exclude_ids=None):
    """选择策略：
    1. 优先选择已有 token 的账号（避免登录）
    2. 遍历队列，找到第一个未被 exclude_ids 包含的账号
    3. 从队列中移除该账号
    4. 返回该账号（由后续逻辑保证最终会重新入队）
    """
    if exclude_ids is None:
        exclude_ids = []

    # 第一轮：优先选择已有 token 的账号
    for i in range(len(account_queue)):
        acc = account_queue[i]
        acc_id = get_account_identifier(acc)
        if acc_id and acc_id not in exclude_ids:
            if acc.get("token", "").strip():  # 已有 token
                logger.info(f"[choose_new_account] 选择已有token的账号: {acc_id}")
                return account_queue.pop(i)

    # 第二轮：选择任意账号（需要登录）
    for i in range(len(account_queue)):
        acc = account_queue[i]
        acc_id = get_account_identifier(acc)
        if acc_id and acc_id not in exclude_ids:
            logger.info(f"[choose_new_account] 选择需登录的账号: {acc_id}")
            return account_queue.pop(i)

    logger.warning("[choose_new_account] 没有可用的账号或所有账号都在使用中")
    return None


def release_account(account: dict):
    """将账号重新加入队列末尾"""
    account_queue.append(account)


# ----------------------------------------------------------------------
# Claude API key 管理函数（简化版本）
# ----------------------------------------------------------------------
def choose_claude_api_key():
    """选择一个可用的Claude API key - 现在直接由用户提供"""
    return None


def release_claude_api_key(api_key):
    """释放Claude API key - 现在无需操作"""
    pass


# ----------------------------------------------------------------------
# 判断调用模式：配置模式 vs 用户自带 token
# ----------------------------------------------------------------------
def determine_mode_and_token(request: Request):
    """
    根据请求头 Authorization 判断使用哪种模式：
    - 如果 Bearer token 出现在 CONFIG["keys"] 中，则为配置模式，从 CONFIG["accounts"] 中随机选择一个账号（排除已尝试账号），
      检查该账号是否已有 token，否则调用登录接口获取；
    - 否则，直接使用请求中的 Bearer 值作为 DeepSeek token。
    结果存入 request.state.deepseek_token；配置模式下同时存入 request.state.account 与 request.state.tried_accounts。
    """
    auth_header = request.headers.get("Authorization", "")
    if not auth_header.startswith("Bearer "):
        raise HTTPException(
            status_code=401, detail="Unauthorized: missing Bearer token."
        )
    caller_key = auth_header.replace("Bearer ", "", 1).strip()
    config_keys = CONFIG.get("keys", [])
    if caller_key in config_keys:
        request.state.use_config_token = True
        request.state.tried_accounts = []  # 初始化已尝试账号
        selected_account = choose_new_account()
        if not selected_account:
            raise HTTPException(
                status_code=429,
                detail="No accounts configured or all accounts are busy.",
            )
        if not selected_account.get("token", "").strip():
            try:
                login_deepseek_via_account(selected_account)
            except Exception as e:
                logger.error(
                    f"[determine_mode_and_token] 账号 {get_account_identifier(selected_account)} 登录失败：{e}"
                )
                raise HTTPException(status_code=500, detail="Account login failed.")

        request.state.deepseek_token = selected_account.get("token")
        request.state.account = selected_account

    else:
        request.state.use_config_token = False
        request.state.deepseek_token = caller_key


def get_auth_headers(request: Request) -> dict:
    """返回 DeepSeek 请求所需的公共请求头"""
    return {**BASE_HEADERS, "authorization": f"Bearer {request.state.deepseek_token}"}


# ----------------------------------------------------------------------
# Claude 认证相关函数
# ----------------------------------------------------------------------
def determine_claude_mode_and_token(request: Request):
    """Claude认证：沿用现有的OpenAI接口认证逻辑"""
    determine_mode_and_token(request)
