const Spinner = (props) => {

    return (
        <div className={props.extraClasses || ""}>
            <img src="/dualRing.svg" />
        </div>
    )
}

export default Spinner