import React, {useContext, useEffect, useState} from "react";
import loadingBarProvider from "./loadingBarProvider";

const LoadingBar = (props) => {

    const [loadingBarStatus, setLoadingBarStatus] = useContext(loadingBarProvider.LoadingBarContext);
    const [progressClass, setProgressClass] = useState("w-0")
    const [displayClass, setDisplayClass] = useState("hidden")

    const classes = [
        "h-[2px]",
        "w-full",
        "bottom-[-2px]",
        "z-[-1]",
        "left-0",
        "absolute",
        props.extraClasses || '',
    ]

    const childClasses = [
        "bg-hospedate-green",
        "h-full",
        "transition-all",
        "duration-300",
        "shadow-[0_0_10px_#227f4b]"
    ]

    useEffect(() => {
        if (loadingBarStatus === loadingBarProvider.LOADING_BAR_LOADING) {
            setDisplayClass("block")
            setTimeout(() => {
                setProgressClass("w-3/4")
            }, 200)
        } else if (loadingBarStatus === loadingBarProvider.LOADING_BAR_READY) {
            setProgressClass("w-full")
            setTimeout(() => {
                setDisplayClass("hidden")
                setProgressClass("w-0")
            }, 200)
        }
    }, [loadingBarStatus])

    return (
        <div className={classes.join(' ')}>
            <div
                className={childClasses.join(' ') + ` ${progressClass} ${displayClass}`}
            >
            </div>
        </div>
    )
}

export default LoadingBar