import Router from 'next/router'
import utils from '../utils/utils'
import ContentBox from '../components/contentBox'
import FormButton from '../components/formButton'
import NavBar from '../components/navBar'
import React, { useContext } from 'react';
import { UserContext, isLoggedIn } from '../components/userProvider'
import FormInputIcon from '../components/formInputIcon'
import toast from "react-hot-toast";

const apiErrors = utils.apiErrors


async function handleSubmit(event, setUser) {
    event.preventDefault()

    const data = {
        email: event.target.email.value,
        password: event.target.password.value,
    }

    const endpoint = '/api/login';
    const options = {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
    }
    const loginResponse = await fetch(endpoint, options)
    const result = await loginResponse.json()

    if (result.error && result.error === apiErrors.EmailOrPasswordIncorrect ) {
        toast.error('Correo electrónico o contraseña incorrectos')
    } else if (result.error) {
        toast.error("Error inesperado")
    } else {
        setUser({
            id: result.data.user.id,
            name: result.data.user.name,
            email: result.data.user.email,
            isHost: result.data.user.isHost,
        });
        Router.push('/')
    }
}


function Login(props) {
    const [user, setUser] = useContext(UserContext);

    if (isLoggedIn(user)) {
        Router.push('/')
    }

    return (
        <div>
            <NavBar />
            <ContentBox>
                <form
                    className="flex flex-col items-center"
                    onSubmit={(event) => handleSubmit(event, setUser)}
                    >
                    <FormInputIcon
                        icon="email"
                        type="email"
                        name="email"
                        id="email"
                        placeholder="Correo electrónico"
                        autofocus={true} />
                    <FormInputIcon
                        icon="password"
                        type="password"
                        name="password"
                        id="password"
                        placeholder="Contraseña" />
                  <div className="mt-8 w-[80%] sm:w-auto flex justify-center">
                    <FormButton type="submit" text="Iniciar sesión" />
                  </div>
                </form>
            </ContentBox>
        </div>
    );
}

export default Login
