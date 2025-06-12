import NavBar from '../components/navBar';
import {useContext, useEffect, useState} from "react";
import {UserContext, isGuest, isLoggedIn} from "../components/userProvider";
import Router from "next/router";
import ContentCard from "../components/contentCard";
import Link from "next/link";
import SmartContainer from "../components/smartContainer";

async function fetchSupportContact() {
    const getPhoneNumber = await fetch(
        "/api/supportContact"
    )
    const response = await getPhoneNumber.json()

    if (response.data) {
        return {
            phoneNumber: response.data.phoneNumber,
            emailAddress: response.data.emailAddress
        }
    } else {
        return null
    }
}

const HelpPage = () => {

    const [user, setUser] = useContext(UserContext);
    const [supportPhoneNumber, setSupportPhoneNumber] = useState("(cargando número)")
    const [supportEmailAddress, setSupportEmailAddress] = useState("(cargando email)")

    if (isGuest(user)) {
        Router.push('/')
    }

    useEffect(() => {
        (async () => {
            const data = await fetchSupportContact()
            if (data) {
                setSupportPhoneNumber(data.phoneNumber)
                setSupportEmailAddress(data.emailAddress)
            }
        })()
    }, [])


    return (
        <div>
            <NavBar />
            <SmartContainer yAxisSpaced={true}>
                {isLoggedIn(user) && (
                    <ContentCard extraClasses="p-6 mt-6 mx-6">
                        <h1 className="text-3xl mb-4">¿Cómo podemos ayudarte?</h1>
                        <p className="pl-4">Queremos que tu experiencia en Hospedate sea la mejor. Por eso, si no encontrás solución a tu problema, podés comunicarte con nosotros mediante correo electrónico a <Link className="text-hosp-light-blue hover:underline cursor-pointer" href={`mailto:${supportEmailAddress}`}>{supportEmailAddress}</Link>, o por WhatsApp a <Link className="text-hosp-light-blue hover:underline cursor-pointer" target="_blank" href={`https://wa.me/${supportPhoneNumber.replace(/ /g, '')}`}>{supportPhoneNumber}</Link>.</p>
                        <h2 className="mt-8 mb-4 text-2xl">Preguntas frecuentes</h2>
                        <p className="pl-4">Visitá nuestra sección de <Link className="text-hosp-light-blue hover:underline cursor-pointer" href="/help/faq">preguntas frecuentes</Link>.</p>
                        <h2 className="mt-8 mb-4 text-2xl">Política de Cancelaciones</h2>
                        <p className="pl-4">Consultá nuestra <Link className="text-hosp-light-blue hover:underline cursor-pointer" href="/help/cancellations">Política de Cancelaciones</Link> para encontrar los reembolsos, retenciones, cargos, compensaciones o modificaciones en la reputación de los usuarios que apliquen debido a cancelaciones de reservas.</p>
                        <h2 className="mt-8 mb-4 text-2xl">Precios de Hospedate</h2>
                        <p className="pl-4">Consultá nuestra página de <Link className="text-hosp-light-blue hover:underline cursor-pointer" href="/help/pricing">Precios</Link>.</p>
                        <h2 className="mt-8 mb-4 text-2xl">Acerca de Hospedate</h2>
                        <p className="pl-4">Hospedate es una plataforma que ofrece un servicio donde los usuarios pueden realizar reservas de alquileres temporarios. Las reservas pueden ser de tipo "directas" o "protegidas". Con las reservas directas ambas partes se pueden contactar directamente y concretar el pago entre ellas. Las reservas protegidas se efectúan a través del sistema de pago con criptomonedas de la plataforma. Para más información podés leer nuestros <Link className="text-hosp-light-blue hover:underline cursor-pointer" href="/help/terms">Términos y Condiciones</Link>.</p>
                    </ContentCard>
                )}
            </SmartContainer>
        </div>
    )
}

export default HelpPage
