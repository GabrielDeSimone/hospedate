const MakeOrderButton = (props) => {

    const handleSubmit = (event, onSubmitCallback) => {
        event.preventDefault();
        if (onSubmitCallback) {
            onSubmitCallback();
        }
    }

    const classes = [
        "border",
        "text-white",
        "font-light",
        "text-lg",
        "w-full",
        "min-h-[45px]",
        "max-w-xs",
        props.bgColor,
        props.extraClasses || '',
    ]

    return (
        <form onSubmit={(event) => handleSubmit(event, props.onSubmit)} className="text-center">
            <button className={classes.join(' ')} type="submit">{props.text}</button>
        </form>
    )
}

export default MakeOrderButton