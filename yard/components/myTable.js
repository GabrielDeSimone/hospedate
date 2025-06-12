const TableHTr = (props) => (
    <tr>{props.children}</tr>
)

const TableHTh = (props) => {
    const classes = [
        "text-xs",
        "text-gray-700",
        "uppercase",
        "bg-gray-50",
        "px-6",
        "py-3",
        props.extraClasses || '',
    ]
    return (
        <th className={classes.join(' ')}>{props.children}</th>
    )
}

const TableHTd = (props) => {
    const classes = ["px-6", "py-4", props.extraClasses || '']
    return (<td className={classes.join(' ')}>{props.children}</td>)
}

const TableH = (props) => {

    const classes = [
        "w-full",
        "text-sm",
        "text-left",
        "text-gray-500",
        "bg-white",
        props.extraClasses || '',
    ]

    return (
        <div className="ml-5 mt-5">
            <table className={classes.join(' ')}>
                <tbody>
                    {props.children}
                </tbody>
            </table>
        </div>
    )
}

const TableV = (props) => {
    const classes = [
        "w-full",
        "text-sm",
        "text-left",
        "text-gray-500",
        props.extraClasses || '',
    ]
    return (<table className={classes.join(' ')}>{props.children}</table>)
}

const TableVHeading = (props) => {
    const classes = [
        "text-xs",
        "text-gray-700",
        "uppercase",
        "bg-gray-50",
    ]
    return (<thead className={classes.join(' ')}><tr>{props.children}</tr></thead>)
}

const TableVThCol = (props) => {
    const classes = ["px-6", "py-3", props.extraClasses || '']
    return (<th scope="col" className={classes.join(' ')}>{props.children}</th>)
}

const TableVTr = (props) => {
    const classes = ["bg-white", "border-b", props.extraClasses || '']
    return (<tr className={classes.join(' ')}>{props.children}</tr>)
}

const TableVThRow = (props) => {
    const classes = [
        "px-6",
        "py-4",
        "font-medium",
        "text-gray-900",
        "whitespace-nowrap",
        props.extraClasses || ''
    ]
    return (<th scope="row" className={classes.join(' ')}>{props.children}</th>)
}

const TableVTd = (props) => {
    const classes = ["px-6", "py-4", props.extraClasses || '']
    return (<td className={classes.join(' ')}>{props.children}</td>)
}

const TableVBody = (props) => (
    <tbody>{props.children}</tbody>
)


export {
    TableH,
    TableHTr,
    TableHTd,
    TableHTh,
    TableV,
    TableVHeading,
    TableVThCol,
    TableVThRow,
    TableVTr,
    TableVTd,
    TableVBody,
}
