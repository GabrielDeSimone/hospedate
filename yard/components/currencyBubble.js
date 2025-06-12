const CurrencyBubble = (props) => {

    const classes = [
        "ring-1",
        "h-[17px]",
        "w-[45px]",
        "rounded-md",
        "text-xs",
        "text-center",
        "text-gray-500",
        "p-[1px]",
        props.extraClasses || '',
    ]

    const currencyColors = {
        "usdt": "ring-[#227F4B]",
    }

    const currencyLabels = {
        "usdt": "USDT",
        "usdc": "USDC",
        "usd": "USD"
    }

    const currencyLabel = currencyLabels[props.currency.toLowerCase()]
    classes.push(currencyColors[props.currency.toLowerCase()])
    return (
        <div className={classes.join(' ')}>
            {currencyLabel}
        </div>
    )
}

export default CurrencyBubble