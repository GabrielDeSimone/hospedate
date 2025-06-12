import CurrencyBubble from "./currencyBubble";

const BalanceLabel = (props) => {

    const classes = [
        "w-[160px]",
        "flex",
        props.extraClasses || '',
    ]

    const currencies = {
        "usdt": "USDT",
        "usdc": "USDC",
        "usd": "USD"
    }

    const currencyLabel = currencies[props.currency.toLowerCase()]

    return (
        <div className={classes.join(' ')}>
            <CurrencyBubble currency={props.currency} extraClasses="mt-3 mr-4" />
            <p className="text-3xl">$ {props.balance}</p>
        </div>
    )
}

export default BalanceLabel