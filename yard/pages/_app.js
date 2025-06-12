import React, { useEffect, useState } from 'react';
import Router from "next/router"
import UserProvider from '../components/userProvider';
import SearchParamsProvider from '../components/searchParamsProvider';
import "../styles/global.css"
import utils from '../utils/utils'
import { Toaster } from 'react-hot-toast';

function App({ Component, pageProps }) {

    const [loadingBarStatus, setLoadingBarStatus] = useState(null);

    Router.events.on("routeChangeStart", () => {
        setLoadingBarStatus(LOADING_BAR_LOADING)
    })

    Router.events.on("routeChangeComplete", () => {
        setLoadingBarStatus(LOADING_BAR_READY)
    })

    // execute after rendered
    useEffect(() => {
        console.log(utils.hospedateStringLogo)
    }, []);

    return (
        <>
            <UserProvider>
                <SearchParamsProvider>
                    <LoadingBarContext.Provider value={[ loadingBarStatus, setLoadingBarStatus ]}>
                        <Component {...pageProps} />
                    </LoadingBarContext.Provider>
                </SearchParamsProvider>
            </UserProvider>
            <Toaster />
        </>
    );
}


export const LOADING_BAR_READY = "ready"
export const LOADING_BAR_LOADING = "loading"

export const LoadingBarContext = React.createContext(null);
export default App;
