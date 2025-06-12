const FormButton = (props) => {

    let classes = [
        "text-white",
        "text-lg",
        "font-light",
        "px-5",
        "py-2.5",
        "mt-4",
        "text-center",
        "disabled:bg-dis-hosp-pale-blue",
        "flex",
        props.extraClasses || '',
    ]

    let bgClasses = [
        "bg-hosp-pale-blue",
    ]

    if (props.smRounded) {
        classes.push("rounded-sm")
    } else {
        classes.push("rounded-lg")
    }

    if (props.bgClasses) {
        bgClasses = props.bgClasses
    }

    return (
        <button type={props.type} disabled={props.disabled} className={classes.concat(bgClasses).join(' ')} onClick={props.onClick}>
            {props.icon ? (
                <span className="
                    material-symbols-outlined
                    text-md
                    bg-transparent
                    h-full
                    mr-3
                ">{props.icon}</span>
            ) : null}
            <span className="block h-full">{props.text}</span>
        </button>
    )
}

export default FormButton