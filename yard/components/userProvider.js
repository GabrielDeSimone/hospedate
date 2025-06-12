import React, { useState, useEffect } from 'react';

async function innerFetchUser() {
    const whoAmIReq = await fetch("/api/whoami")
    const response = await whoAmIReq.json()

    const user = response.data.user
    return user ? {
        id: user.id,
        name: user.name,
        email: user.email,
        isHost: user.isHost,
    } : {
        id: 'guest'
    }
}

const UserProvider = (props) => {
    const [user, setUser] = useState(null);

    useEffect(() => {
        async function fetchUser(){
            const user = await innerFetchUser();
            setUser(user);
        }
        fetchUser()
    }, []);

    return (
        <UserContext.Provider value={[ user, setUser ]}>
            {props.children}
        </UserContext.Provider>
    );
}

export const isGuest = (user) => Boolean(user && user.id === 'guest')
export const isLoggedIn = (user) => Boolean(user && user.id !== 'guest')

export const isHost = (user) => Boolean(isLoggedIn(user) && user.isHost)

export const isNotHost = (user) => Boolean(isLoggedIn(user) && !user.isHost)

export const UserContext = React.createContext(null);
export default UserProvider