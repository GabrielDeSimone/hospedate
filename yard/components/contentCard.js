const ContentCard = (props) => {

    let classes = [
        "border",
        "border-gray-200",
        'rounded-lg',
        'p-6',
        props.extraClasses || '',
    ]

    if (props.shadow) {
        classes = [
            ...classes,
            'shadow-lg',
        ]
    }

    return (
        <div className={classes.join(' ')}>
            {props.children}
        </div>
    )
}

export default ContentCard