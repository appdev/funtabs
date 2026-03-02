import React, { memo, useCallback, useEffect, useMemo, useRef, useState } from "react";
import { Flipper } from "react-flip-toolkit";
import ReactGridLayout from "react-grid-layout/legacy";
import "react-grid-layout/css/styles.css";
import { useSelector } from "react-redux";
import Slider from "react-slick";
import "slick-carousel/slick/slick-theme.css";
import "slick-carousel/slick/slick.css";
import "react-resizable/css/styles.css";
import "../common/funtabs.css";
import DefaultStyle from "../showStyle/defaultStyle";
import OnlyIconStyle from "../showStyle/onlyIconStyle";
import OnlyText from "../showStyle/onlyText";
import PhoneStyle from "../showStyle/phoneStyle";

// 根据卡片类型和显示样式，返回对应的格子宽高（w=列数, h=行数）
// size 字符串格式：第1位=行数(h)，第2位=列数(w)，例如 "32" = 3行2列
function getCardWH(item, cardStyle) {
	const sizeTable = {
		defaultCard: {
			note: "32", hotEventList: "32", translatelite: "32",
			timeProgress: "21", localWeather: "21", festivalRemind: "21", markdown: "11",
		},
		onlyIconCard: {
			note: "34", hotEventList: "34", translatelite: "34",
			timeProgress: "22", localWeather: "22", festivalRemind: "22", markdown: "11",
		},
		phoneCard: {
			note: "34", hotEventList: "34", translatelite: "34",
			timeProgress: "22", localWeather: "22", festivalRemind: "22", markdown: "11",
		},
		onlyText: {
			timeProgress: "11", localWeather: "11", markdown: "11",
		},
	};
	const table = sizeTable[cardStyle] ?? sizeTable.defaultCard;
	const code = String(table[item.type] ?? item.size ?? "11");
	return { h: parseInt(code[0]) || 1, w: parseInt(code[1]) || 1 };
}

// 生成 react-grid-layout 所需的布局数组
// 有保存整数坐标的卡片优先放到指定位置，其余按数组顺序自动填充空格
function buildLayout(cards, cols, cardStyle) {
	if (!cols || cols <= 0) return [];

	const occupied = new Set();

	function isOccupied(x, y, w, h) {
		for (let r = y; r < y + h; r++) {
			for (let c = x; c < x + w; c++) {
				if (c >= cols || occupied.has(`${r},${c}`)) return true;
			}
		}
		return false;
	}

	function occupy(x, y, w, h) {
		for (let r = y; r < y + h; r++) {
			for (let c = x; c < x + w; c++) {
				occupied.add(`${r},${c}`);
			}
		}
	}

	function findFirstFit(w, h) {
		for (let row = 0; row < 9999; row++) {
			for (let col = 0; col <= cols - w; col++) {
				if (!isOccupied(col, row, w, h)) return { x: col, y: row };
			}
		}
		return { x: 0, y: 0 };
	}

	const fixed = [];
	const auto = [];

	cards.forEach((card) => {
		const { w, h } = getCardWH(card, cardStyle);
		const sx = typeof card.x === "number" ? card.x : parseInt(card.x);
		const sy = typeof card.y === "number" ? card.y : parseInt(card.y);
		const withinBounds =
			Number.isInteger(sx) && Number.isInteger(sy) &&
			sx >= 0 && sy >= 0 && sx + w <= cols;
		if (withinBounds) {
			fixed.push({ card, sx, sy, w, h });
		} else {
			auto.push({ card, w, h });
		}
	});

	const layout = [];

	// 先放置有保存坐标的卡片，如果位置已被占用则退回到自动排列
	fixed.forEach(({ card, sx, sy, w, h }) => {
		if (!isOccupied(sx, sy, w, h)) {
			layout.push({ i: String(card.id), x: sx, y: sy, w, h });
			occupy(sx, sy, w, h);
		} else {
			auto.push({ card, w, h });
		}
	});

	// 自动排列其余卡片，按数组顺序填入第一个空位
	auto.forEach(({ card, w, h }) => {
		const { x, y } = findFirstFit(w, h);
		layout.push({ i: String(card.id), x, y, w, h });
		occupy(x, y, w, h);
	});

	return layout;
}

const ShowList = memo((props) => {
	const {
		cardStyle,
		tabs,
		gap,
		num,
		newlinkList,
		setNewLinkList,
		drag,
		edit,
		radius,
		widthNum,
		heightNum,
		gridWidthNum,
		tabsActiveKey,
		setTabsActiveKey,
		linkOpen,
		fontColor,
		editFunction,
	} = props;

	const deviceType = useSelector((state) => state.deviceType.type);
	const [a, setA] = useState(0);
	const [folderDisabled, setFolderDisabled] = useState(true);
	const [click] = useState(0);
	const containerRef = useRef(null);
	const [containerWidth, setContainerWidth] = useState(0);

	useEffect(() => {
		setFolderDisabled(edit !== "");
	}, [edit]);

	useEffect(() => {
		setTimeout(() => {
			setA(1);
		}, 500);
	}, []);

	// 监听容器宽度变化，用于自动计算列数
	useEffect(() => {
		if (!containerRef.current) return;
		const observer = new ResizeObserver((entries) => {
			setContainerWidth(entries[0].contentRect.width);
		});
		observer.observe(containerRef.current);
		return () => observer.disconnect();
	}, []);

	function showKey(val) {
		return deviceType === "PC" && val !== 0 ? newlinkList : false;
	}

	const settings = {
		dots: false,
		infinite: false,
		adaptiveHeight: true,
		touchMove: false,
		initialSlide: num,
		arrows: false,
	};

	// 根据可用宽度和卡片宽度计算能放几列（等同于 CSS auto-fill 效果）
	const cols = useMemo(() => {
		const gridPadding = 16; // #grid-div 左右各 8px padding
		const available = containerWidth > gridPadding ? containerWidth - gridPadding : containerWidth;
		const maxWidth = gridWidthNum ? Math.min(available, parseInt(gridWidthNum)) : available;
		return Math.max(1, Math.floor((maxWidth + gap) / (widthNum + gap)));
	}, [containerWidth, gap, widthNum, gridWidthNum]);

	// 网格实际宽度：列数 × 卡片宽 + (列数-1) × 间距
	const gridWidth = useMemo(() => {
		return cols * widthNum + Math.max(0, cols - 1) * gap;
	}, [cols, widthNum, gap]);

	// 拖拽结束后把新坐标回写到卡片数据，下次渲染时保持位置
	const handleLayoutChange = useCallback(
		(layout, tabIndex) => {
			setNewLinkList((prev) => {
				const updated = prev.map((tab, i) => {
					if (i !== tabIndex) return tab;
					const newContent = tab.content.map((card) => {
						const item = layout.find((l) => l.i === String(card.id));
						if (!item) return card;
						return { ...card, x: item.x, y: item.y };
					});
					return { ...tab, content: newContent };
				});
				return updated;
			});
		},
		[setNewLinkList]
	);

	function renderCard(item, index, tabIndex) {
		const commonProps = {
			id: index,
			edit,
			item,
			linkList: newlinkList,
			setLinkList: setNewLinkList,
			radius,
			num: tabIndex,
			gap,
			widthNum,
			heightNum,
			tabsActiveKey,
			setTabsActiveKey,
			cardStyle,
			linkOpen,
			setA,
			folderDisabled,
			click,
			editFunction,
		};

		switch (cardStyle) {
			case "defaultCard":
				return <DefaultStyle key={item.link + item.type + item.id} {...commonProps} />;
			case "onlyIconCard":
				return <OnlyIconStyle key={item.link + item.type + item.id} {...commonProps} />;
			case "phoneCard":
				return <PhoneStyle key={item.link + item.type + item.id} {...commonProps} />;
			case "onlyText":
				return <OnlyText key={item.link + item.type + item.id} {...commonProps} fontColor={fontColor} />;
			default:
				return null;
		}
	}

	function creatLinkList(tab, tabIndex) {
		const layout = buildLayout(tab.content, cols, cardStyle);
		return (
			<div key={tabIndex}>
				<div
					id="grid-div"
					style={{
						display: "flex",
						justifyContent: "center",
						padding: "8px",
					}}
				>
					<div style={{ width: gridWidth || widthNum, flexShrink: 0, pointerEvents: "all" }}>
						<ReactGridLayout
							layout={layout}
							cols={cols}
							rowHeight={heightNum}
							width={gridWidth || widthNum}
							margin={[gap, gap]}
							containerPadding={[0, 0]}
							isDraggable={!drag}
							isResizable={false}
							compactType={null}
							preventCollision={true}
							onDragStop={(layout) => handleLayoutChange(layout, tabIndex)}
							style={{ minHeight: `${heightNum}px` }}
						>
							{tab.content.map((item, index) => (
								<div key={String(item.id)} style={{ width: "100%", height: "100%" }}>
									{renderCard(item, index, tabIndex)}
								</div>
							))}
						</ReactGridLayout>
					</div>
				</div>
			</div>
		);
	}

	return (
		<>
			{/* 零高度占位 div，用于测量可用容器宽度 */}
			<div ref={containerRef} style={{ width: "100%", height: 0, pointerEvents: "none" }} />
			<Flipper flipKey={showKey(a)} spring="veryGentle">
				<Slider ref={tabs} {...settings}>
					{newlinkList.map(creatLinkList)}
				</Slider>
			</Flipper>
		</>
	);
});

export default ShowList;
