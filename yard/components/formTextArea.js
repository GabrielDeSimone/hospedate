import React from 'react'

const FormTextArea = (props) => {

    let widthClasses = "w-full sm:w-[80%]"
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
        props.extraClasses ? props.extraClasses : ""
    ]

    if (props.widthClasses) {
        widthClasses = props.widthClasses
    }

    return (
        <div className={`flex flex-row mt-4 ${widthClasses}`}>
            <label htmlFor={props.id} className="cursor-pointer">
                <span className="
                    material-symbols-outlined
                    text-3xl
                    bg-transparent
                    text-gray-700
                ">{props.icon}</span>
            </label>
            <textarea
                name={props.name}
                id={props.id}
                onChange={props.onChange}
                required={props.required}
                value={props.value}
                className={classes.join(' ')}
            />
        </div>
    )
}

export default FormTextArea