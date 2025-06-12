const SmartContainer = (props) => {

    let classes = [
        "laptop:max-w-[1080px]",
        "desktop:max-w-[1392px]",
        "w-auto",
        "mx-auto",
        "px-6",
        props.extraClasses || '',
    ]

    if (props.roundedBorder) {
        classes = [
            ...classes,
            "rounded-lg",
            "border",
            "border-gray-300",
        ]
    }

    if (props.yAxisSpaced) {
        classes = [
            ...classes,
            "mt-6",
        ]
    }

    if (! props.noTabletWidthLimit) {
        classes = [
            ...classes,
            "tablet:max-w-[696px]",
        ]
    }

    return (
        <div className={classes.join(' ')}>
            {props.children}
        </div>
    )
}

export default SmartContainer