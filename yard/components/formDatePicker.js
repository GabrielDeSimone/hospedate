import Datepicker from "tailwind-datepicker-react"

/**
 * 
 * Not being used for now. To try it out again, please install tailwind-datepicker-react
 */

const FormDatePicker = (props) => {
    const options = {
        title: props.title,
        autoHide: true,
        todayBtn: false,
        clearBtn: false,
        maxDate: new Date("2025-01-01"),
        minDate: new Date(),
        theme: {
            background: "",
            todayBtn: "",
            clearBtn: "",
            icons: "",
            text: "",
            disabledText: "",
            input: "",
            inputIcon: "",
            selected: "",
        },
        datepickerClassNames: "top-12 mb-6",
        defaultDate: new Date(),
        language: "es",
    }


    return (
        <div className="mb-6">
            <Datepicker options={options} onChange={props.onChange} show={props.show} setShow={props.setShow} />
        </div>
    )
}

export default FormDatePicker