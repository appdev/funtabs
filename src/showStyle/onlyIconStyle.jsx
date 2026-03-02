import { memo } from "react";
import { Flipped } from "react-flip-toolkit";
import { useInView } from "react-intersection-observer";
import { useLongPress } from "use-long-press";
import "../common/funtabs.css";
import Article from "../component/article";
import FestivalRemind from "../component/festivalRemind";
import HotEventList from "../component/hotEventList";
import LocalWeather from "../component/localWeather";
import Markdown from "../component/markdown";
import Note from "../component/note";
import TimeProgress from "../component/timeProgress";
import Translate from "../component/translate";
import LinkCard from "../module/linkCard";

const OnlyIconStyle = memo((props) => {
	const {
		id,
		edit,
		item,
		num,
		gap,
		linkList,
		setLinkList,
		radius,
		widthNum,
		heightNum,
		cardStyle,
		linkOpen,
		click,
		setDisabled,
		setA,
		folderDisabled,
		type,
		editFunction,
		tabsActiveKey,
		setTabsActiveKey,
	} = props;

	const longpress = useLongPress(
		() => {
			if (edit === "none") {
				editFunction();
			}
		},
		{
			cancelOnMovement: true,
			detect: "mouse",
		}
	);

	const [ref, inView] = useInView({
		threshold: 0,
		rootMargin: `${window.innerHeight / 2}px`,
	});

	function renderComponent(Component, className) {
		return (
			<div
				ref={ref}
				style={{
					width: "100%",
					height: "100%",
					position: "relative",
					pointerEvents: "all",
				}}
				{...longpress()}
			>
				{inView && (
					<Component
						id={id}
						edit={edit}
						item={item}
						num={num}
						linkList={linkList}
						setLinkList={setLinkList}
						radius={radius}
						heightNum={heightNum}
						cardStyle={cardStyle}
						linkOpen={linkOpen}
						click={click}
						setDisabled={setDisabled}
						widthNum={widthNum}
						gap={gap}
						setClick={setA}
						folderDisabled={folderDisabled}
						type={type}
						tabsActiveKey={tabsActiveKey}
						setTabsActiveKey={setTabsActiveKey}
					/>
				)}
			</div>
		);
	}

	function howToShow() {
		switch (item.type) {
			case "link":
				return renderComponent(LinkCard, "");
			case "note":
				return renderComponent(Note, "");
			case "timeProgress":
				return renderComponent(TimeProgress, "");
			case "markdown":
				return renderComponent(Markdown, "");
			case "translatelite":
				return renderComponent(Translate, "");
			case "localWeather":
				return renderComponent(LocalWeather, "");
			case "hotEventList":
				return renderComponent(HotEventList, "");
			case "festivalRemind":
				return renderComponent(FestivalRemind, "");
			case "article":
				return renderComponent(Article, "");
			default:
				return null;
		}
	}

	return (
		<>
			<Flipped flipId={item.type + item.link + item.id}>{howToShow()}</Flipped>
		</>
	);
});

export default OnlyIconStyle;
