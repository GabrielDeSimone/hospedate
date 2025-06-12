import Router from 'next/router'
import ContentBox from '../components/contentBox'
import utils from '../utils/utils'
import FormButton from '../components/formButton'
import NavBar from '../components/navBar'
import FormInputIcon from '../components/formInputIcon'
import { useContext } from 'react'
import { UserContext, isLoggedIn } from '../components/userProvider'
import toast from "react-hot-toast";

const apiErrors = utils.apiErrors

function Register() {
    const [user, setUser] = useContext(UserContext);

    if (isLoggedIn(user)) {
        Router.push('/')
    }

    const handleSubmit = async (event) => {
        event.preventDefault()

        // check if repeat password is correct
        if (event.target.password.value !== event.target.repeat_password.value) {
            toast.error('Las contraseñas no coinciden')
            return
        }

        const data = {
            email: event.target.email.value,
            name: event.target.name.value,
            password: event.target.password.value,
            phoneNumber: event.target.phone_number.value,
            invitationId: event.target.invitation_id.value,
        }

        const JSONdata = JSON.stringify(data)
        const endpoint = '/api/register'

        const options = {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSONdata,
        }
        const response = await fetch(endpoint, options)
        const result = await response.json()

        if (result.error && result.error === apiErrors.EmailOrPhoneAlreadyExist) {
            toast.error('El correo electrónico o el celular ya están registrados')
        } else if (result.error && result.error === apiErrors.InvitationNotValid) {
            toast.error('La invitación ingresada no es válida')
        } else if (result.error) {
            toast.error("Ha ocurrido un error inesperado. Intente más tarde")
        } else {
            toast.success("Registro exitoso")
            Router.push('/login')
        }
    }

    return (
        <div>
          <NavBar />
          <ContentBox>
                <form className="flex flex-col items-center" onSubmit={handleSubmit}>
                  <FormInputIcon icon="email" type="email" name="email" id="email" placeholder="Correo electrónico" required />
                  <FormInputIcon icon="password" type="password" name="password" id="password" placeholder="Contraseña" minLength="8" required />
                  <FormInputIcon icon="password" type="password" name="repeat_password" id="repeat_password" placeholder="Repetir contraseña" minLength="8" required />
                  <FormInputIcon icon="badge" name="name" id="name" placeholder="Nombre" required />
                  <FormInputIcon icon="smartphone" name="phone_number" id="phone_number" placeholder="Celular" required />
                  <FormInputIcon icon="partner_exchange" name="invitation_id" id="invitation_id" placeholder="Invitación" required />
                  <div className="mt-8 w-[80%] sm:w-auto flex justify-center">
                    <FormButton type="submit" text="Registrarse" />
                  </div>
                </form>
          </ContentBox>
        </div>
    );
}

export default Register