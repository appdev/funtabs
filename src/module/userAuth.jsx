import {
	ArrowLeftOutlined,
	CloudOutlined,
	UserOutlined,
} from "@ant-design/icons";
import {
	Avatar,
	Button,
	Dropdown,
	Input,
	Modal,
	Radio,
	Space,
	Switch,
	message,
} from "antd";
import React, { useCallback, useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { get, post } from "../common/fetch";
import getLocalData from "../common/getLocalData";
import variable from "../common/variable";
import { setAutoCloudSync } from "../redux/slice/cloudSync";

const VIEW_LOGIN = "login";
const VIEW_REGISTER = "register";

function UserAuth({ fontColor }) {
	const dispatch = useDispatch();
	const autoCloudSync = useSelector((state) => state.cloudSync.autoCloudSync);

	const [modalOpen, setModalOpen] = useState(false);
	const [view, setView] = useState(VIEW_LOGIN);
	const [loading, setLoading] = useState(false);
	const [token, setToken] = useState(() => localStorage.getItem("token"));
	const [userInfo, setUserInfo] = useState(null);

	// 登录表单
	const [username, setUsername] = useState("");
	const [password, setPassword] = useState("");
	const [dataOption, setDataOption] = useState("pull");

	// 注册表单
	const [regUsername, setRegUsername] = useState("");
	const [regPassword, setRegPassword] = useState("");

	// 拉取用户信息
	const loadUserInfo = useCallback((tk) => {
		get(variable.getUserInfo, null, tk)
			.then((res) => {
				if (res) setUserInfo(res);
			})
			.catch(() => {});
	}, []);

	useEffect(() => {
		if (token) loadUserInfo(token);
	}, [token, loadUserInfo]);

	// 将本地数据推送到云端
	async function pushToCloud(tk) {
		const data = await getLocalData();
		await post(variable.saveData, { data: JSON.stringify(data) }, tk);
	}

	// 从云端拉取数据并覆盖本地
	async function pullFromCloud(tk) {
		const res = await get(variable.getData, null, tk);
		if (!res) return false;
		const data = typeof res === "string" ? JSON.parse(res) : res;
		if (typeof data === "object") {
			Object.entries(data).forEach(([key, value]) => {
				if (key !== "token" && key !== "activeKey") {
					localStorage.setItem(key, value);
				}
			});
		}
		return true;
	}

	// 登录
	async function handleLogin() {
		if (!username || !password) {
			message.warning("请填写用户名和密码");
			return;
		}
		setLoading(true);
		try {
			const formData = new FormData();
			formData.append("username", username);
			formData.append("password", password);
			const res = await post(variable.login, formData);
			if (!res) return;

			const tk = res.token || res;
			if (!tk || typeof tk !== "string") {
				message.error("登录失败，请重试");
				return;
			}

			localStorage.setItem("token", tk);
			setToken(tk);
			loadUserInfo(tk);

			if (dataOption === "push") {
				await pushToCloud(tk);
				message.success("登录成功，本地数据已上传云端");
			} else if (dataOption === "pull") {
				const ok = await pullFromCloud(tk);
				if (ok) {
					message.success("登录成功，已获取云端数据");
					setTimeout(() => window.location.reload(), 800);
				} else {
					message.success("登录成功");
				}
			} else {
				message.success("登录成功");
			}

			dispatch(setAutoCloudSync(true));
			setModalOpen(false);
		} finally {
			setLoading(false);
		}
	}

	// 注册（仅用户名+密码）
	async function handleRegister() {
		if (!regUsername || !regPassword) {
			message.warning("请填写用户名和密码");
			return;
		}
		setLoading(true);
		try {
			const formData = new FormData();
			formData.append("username", regUsername);
			formData.append("password", regPassword);
			await post(variable.register, formData);
			message.success("注册成功，请登录");
			setView(VIEW_LOGIN);
		} finally {
			setLoading(false);
		}
	}

	// 退出登录
	function handleLogout() {
		localStorage.removeItem("token");
		setToken(null);
		setUserInfo(null);
		dispatch(setAutoCloudSync(false));
		message.success("已退出登录");
	}

	function openModal() {
		setView(VIEW_LOGIN);
		setModalOpen(true);
	}

	// ---------- 渲染：登录表单 ----------
	function renderLogin() {
		return (
			<>
				<div style={{ textAlign: "center", marginBottom: 20 }}>
					<img
						src="/logo.svg"
						alt="logo"
						style={{ width: 56, height: 56 }}
						onError={(e) => {
							e.target.style.display = "none";
						}}
					/>
					<div
						style={{
							color: "#1677ff",
							fontWeight: "bold",
							marginTop: 8,
							fontSize: 15,
						}}
					>
						山海天地间，与君终相见！
					</div>
				</div>
				<Input
					placeholder="用户名"
					value={username}
					onChange={(e) => setUsername(e.target.value)}
					style={{ marginBottom: 12 }}
					autoComplete="off"
				/>
				<Input.Password
					placeholder="密码"
					value={password}
					onChange={(e) => setPassword(e.target.value)}
					style={{ marginBottom: 14 }}
					onPressEnter={handleLogin}
				/>
				<div style={{ marginBottom: 16, fontSize: 13 }}>
					<span>数据选项：</span>
					<Radio.Group
						value={dataOption}
						onChange={(e) => setDataOption(e.target.value)}
						size="small"
					>
						<Radio value="push">本地上传</Radio>
						<Radio value="pull">获取云端</Radio>
						<Radio value="skip">暂不同步</Radio>
					</Radio.Group>
				</div>
				<Button
					type="primary"
					block
					loading={loading}
					onClick={handleLogin}
					style={{ marginBottom: 14 }}
				>
					登 录
				</Button>
				<div style={{ textAlign: "right" }}>
					<Button
						type="link"
						size="small"
						onClick={() => setView(VIEW_REGISTER)}
					>
						新用户注册
					</Button>
				</div>
			</>
		);
	}

	// ---------- 渲染：注册表单 ----------
	function renderRegister() {
		return (
			<>
				<div
					style={{
						color: "#1677ff",
						fontWeight: "bold",
						textAlign: "center",
						fontSize: 15,
						marginBottom: 20,
					}}
				>
					欢迎您的注册使用
				</div>
				<Input
					placeholder="用户名"
					value={regUsername}
					onChange={(e) => setRegUsername(e.target.value)}
					style={{ marginBottom: 12 }}
					autoComplete="off"
				/>
				<Input.Password
					placeholder="密码"
					value={regPassword}
					onChange={(e) => setRegPassword(e.target.value)}
					style={{ marginBottom: 20 }}
					onPressEnter={handleRegister}
				/>
				<Button
					type="primary"
					block
					loading={loading}
					onClick={handleRegister}
				>
					注 册
				</Button>
			</>
		);
	}

	// ---------- 已登录：下拉菜单 ----------
	const userMenuItems = [
		{
			key: "sync",
			label: (
				<div
					style={{
						display: "flex",
						alignItems: "center",
						justifyContent: "space-between",
						gap: 24,
						padding: "2px 0",
					}}
				>
					<span>云端同步</span>
					<Switch
						size="small"
						checked={autoCloudSync}
						onChange={(checked) => dispatch(setAutoCloudSync(checked))}
						onClick={(_, e) => e.stopPropagation()}
					/>
				</div>
			),
		},
		{ type: "divider" },
		{
			key: "logout",
			label: "退出登录",
			danger: true,
		},
	];

	// 已登录态：头像 + 用户名
	if (token) {
		return (
			<Dropdown
				menu={{
					items: userMenuItems,
					onClick: ({ key }) => {
						if (key === "logout") handleLogout();
					},
				}}
				placement="bottomRight"
				trigger={["click"]}
			>
				<Button
					type="text"
					style={{
						color: fontColor,
						padding: "0 8px",
						display: "flex",
						alignItems: "center",
						gap: 6,
						marginRight: "-10px",
					}}
				>
					<Avatar
						size={22}
						src={userInfo?.avatar}
						icon={<UserOutlined />}
						style={{ backgroundColor: "#1677ff", flexShrink: 0 }}
					/>
					<span
						style={{
							fontSize: 13,
							maxWidth: 64,
							overflow: "hidden",
							textOverflow: "ellipsis",
							whiteSpace: "nowrap",
						}}
					>
						{userInfo?.username || "用户"}
					</span>
				</Button>
			</Dropdown>
		);
	}

	// 未登录态：登录按钮 + 弹窗
	const modalTitle =
		view === VIEW_REGISTER ? (
			<Space>
				<Button
					type="text"
					icon={<ArrowLeftOutlined />}
					size="small"
					onClick={() => setView(VIEW_LOGIN)}
				/>
				用户注册
			</Space>
		) : null;

	return (
		<>
			<Button
				type="text"
				style={{
					color: fontColor,
					display: "flex",
					alignItems: "center",
					gap: 4,
					fontWeight: "bold",
					marginRight: "-10px",
				}}
				onClick={openModal}
			>
				<CloudOutlined />
				登录 / 注册
			</Button>
			<Modal
				open={modalOpen}
				onCancel={() => setModalOpen(false)}
				footer={null}
				width={380}
				destroyOnClose
				title={modalTitle}
				styles={{ body: { paddingTop: view === VIEW_LOGIN ? 8 : 4 } }}
			>
				{view === VIEW_LOGIN && renderLogin()}
				{view === VIEW_REGISTER && renderRegister()}
			</Modal>
		</>
	);
}

export default UserAuth;
