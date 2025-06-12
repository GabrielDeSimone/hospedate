import NavBar from '../../components/navBar';
import {useContext, useState} from "react";
import {UserContext, isGuest, isLoggedIn} from "../../components/userProvider";
import Router from "next/router";
import ContentCard from "../../components/contentCard";
import Link from "next/link";
import SmartContainer from "../../components/smartContainer";

const PricingPage = () => {

    const [user, setUser] = useContext(UserContext);

    if (isGuest(user)) {
        Router.push('/')
    }

    return (
        <div>
            <NavBar />
            <SmartContainer yAxisSpaced={true}>
                {isLoggedIn(user) && (
                    <ContentCard extraClasses="helpTextSection">
                        <Link className="text-hosp-light-blue hover:underline cursor-pointer mb-4 block" href="/help">← volver a la página principal de ayuda</Link>
                        <h1>Precios de Hospedate</h1>
                        <p>Nuestros cargos se clasifican en dos tipos. Cargos por uso del servicio (comisiones en reservas protegidas) o cargos por usar Hospedate en modo anfitrión.</p>
                        <h2>Cargos por el uso del servicio</h2>
                        <p>En Hospedate existen dos tipos de reservas: directas y protegidas. Las reservas directas están exentas de todo tipo de comisión. Los cargos por el uso del servicio aplican únicamente para reservas protegidas. La comisión para reservas protegidas está dividida en Anfitrión y Huésped, al igual que en otras plataformas.</p>
                        <h3>Comisión para anfitriones</h3>
                        <p>La comisión para anfitriones es del 1.5% y es deducida del subtotal. El subtotal se define como el precio por noche por la cantidad de noches. Esto determina la ganancia de los anfitriones, salvo impuestos que deban pagar de acuerdo a su jurisdicción.</p>
                        <h3>Comisión para huéspedes</h3>
                        <p>La comisión para huéspedes es del 7% y es aplicada sobre el subtotal. El subtotal se define como el precio por noche por la cantidad de noches. Esto determina el costo a pagar por los huéspedes. Los huéspedes deben tener en cuenta las comisiones de la red de criptomonedas a la hora de hacer una transacción, que en el caso de la red TRON (red con la que trabaja Hospedate) es aproximadamente de USDT $ 1.</p>
                        <h2>Cargos por usar Hospedate en modo anfitrión</h2>
                        <p>El uso de Hospedate en modo anfitrión tiene un costo de USD $3 por mes. Sin embargo, todos los anfitriones que ingresen a la plataforma durante el periodo de prueba tienen una bonificación del 100% durante 12 meses.</p>
                    </ContentCard>
                )}
            </SmartContainer>
        </div>
    )
}

export default PricingPage
