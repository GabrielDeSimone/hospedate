import React from 'react'

const FormSelect = (props) => {

    const classes = [
        "w-full",
        "border rounded-lg",
        "pl-3",
        "ml-3",
        "focus:ring-1",
        "focus:ring-hospedate-green",
        "focus:border-transparent",
        "hover:ring-hospedate-green",
        "hover:ring-1",
        "focus:outline-none",
        "text-gray-400",
        "focus:text-gray-700",
        props.extraClasses ? props.extraClasses : " "
    ]

    let widthClasses = "w-[80%] sm:w-[300px]"
    if (props.widthClasses) {
        widthClasses = props.widthClasses
    }

    const containerClasses = [
        "flex",
        "flex-row",
        "mt-4",
        widthClasses,
        props.containerExtraClasses ? props.containerExtraClasses : " "
    ]

    return (
        <div className={containerClasses.join(' ')}>
            <label htmlFor={props.id} className="cursor-pointer">
                <span className="
                    material-symbols-outlined
                    text-3xl
                    bg-transparent
                    text-gray-700
                ">{props.icon}</span>
            </label>
            <select
                name={props.name}
                id={props.id}
                onChange={props.onChange}
                required={props.required}
                value={props.value}
                className={classes.join(' ')}
            >
                {props.children}
            </select>
        </div>
    )
}

export default FormSelect