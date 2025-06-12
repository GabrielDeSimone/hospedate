import React, { useEffect } from 'react'

const FormInputIcon = (props) => {

    const inputType = props.type || "text"
    let datalist = null;
    const datalistId = "datalistId"

    let widthClasses = "w-[80%] sm:w-[300px]"
    if (props.widthClasses) {
        widthClasses = props.widthClasses
    }

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

    const containerClasses = [
        "flex",
        "flex-row",
        "mt-4",
        widthClasses,
        props.containerExtraClasses ? props.containerExtraClasses : " "
    ]

    useEffect(() => {
        if (props.autofocus) {
            document.getElementById(props.id).focus()
        }
    }, [])

    if (inputType === "datalist") {
        const datalistOptions = props.datalistOptions;
        datalist = (<datalist id={datalistId}>
            {datalistOptions.map(option => (<option key={option} value={option} />))}
        </datalist>)
    }



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
            <input
                type={inputType}
                name={props.name}
                id={props.id}
                placeholder={props.placeholder}
                autoFocus={props.autofocus}
                required={props.required}
                value={props.value}
                onChange={props.onChange}
                list={datalist ? datalistId : null}
                minLength={props.minLength ? props.minLength : null}
                className={classes.join(' ')}
            />
            {datalist}
        </div>
    )
}

export default FormInputIcon