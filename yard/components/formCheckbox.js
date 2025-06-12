import React from 'react'

const FormCheckbox = (props) => {

    const classes = [
        "border rounded-lg",
        "pl-3",
        "ml-3",
        "mt-2",
        "focus:ring-1",
        "focus:ring-hospedate-green",
        "focus:border-transparent",
        "hover:ring-hospedate-green",
        "hover:ring-1",
        "focus:outline-none",
        "text-gray-400",
        "focus:text-gray-700",
        "block",
        "h-[20px]",
        "w-[20px]",
        props.extraClasses ? props.extraClasses : " "
    ]

    return (
        <div className="flex flex-row mt-4 w-full sm:w-auto">
            <label htmlFor={props.id} className="cursor-pointer flex">
                <span className="
                    material-symbols-outlined
                    text-3xl
                    bg-transparent
                    text-gray-700
                ">{props.icon}</span>
                <input
                    type="checkbox"
                    checked={props.checked}
                    onChange={props.onChange}
                    className={classes.join(' ')}
                />
                <span className="ml-3 mt-2 block">{props.label}</span>
            </label>
        </div>
    )
}

export default FormCheckbox