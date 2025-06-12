import NavBar from '../../components/navBar';
import {useContext, useState} from "react";
import {UserContext, isGuest, isLoggedIn} from "../../components/userProvider";
import Router from "next/router";
import ContentCard from "../../components/contentCard";
import Link from "next/link";
import SmartContainer from "../../components/smartContainer";

const CancellationsPage = () => {

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
                        <h1 className="text-3xl mb-4">Política de Cancelaciones para reservas</h1>
                        <p>Esta política define los reembolsos, retenciones, cargos, compensaciones o modificaciones en la reputación de los usuarios en consecuencia a cancelaciones de reservas.</p>
                        <p>Se establecen tres rangos en función de la anticipación con que se efectúa la cancelación.</p>
                        <ol>
                            <li>Rango Verde: Cuando la cancelación se efectúa 15 días antes de la fecha de check-in o antes.</li>
                            <li>Rango Amarillo: Cuando la cancelación se efectúa 14 días antes de la fecha de check-in hasta 48 horas antes.</li>
                            <li>Rango Rojo: Cuando la cancelación se efectúa dentro de las 48 horas antes de la fecha de check-in.</li>
                        </ol>
                        <p>De acuerdo a los siguientes rangos temporales, las consecuencias de una cancelación pueden variar para huéspedes o anfitriones.</p>
                        <h2>Política de cancelación para reservas directas</h2>
                        <ol>
                            <li>Rango Verde: Cancelación sin cargos.</li>
                            <li>Rango Amarillo: Cancelación sin cargos.</li>
                            <li>Rango Rojo: Se aplicará una modificación de la reputación del usuario que haya efectuado la cancelación.</li>
                        </ol>
                        <h2>Política de cancelación para reservas protegidas</h2>
                        <h3>Cancelaciones por parte de huéspedes</h3>
                        <ol>
                            <li>Rango Verde: Se aplicará un reembolso total del monto de la reserva para el huésped, menos comisiones.</li>
                            <li>Rango Amarillo:</li>
                            <ul>
                                <li>Se aplicará un reembolso del 70% del monto de la reserva menos comisiones.</li>
                                <li>Se aplicará una compensación para el anfitrión del 30% del monto de la reserva menos comisiones.</li>
                            </ul>
                            <li>Rango Rojo</li>
                            <ul>
                                <li>Se aplicará un reembolso del 50% del monto de la reserva menos comisiones.</li>
                                <li>Se aplicará una compensación para el anfitrión del 50% del monto de la reserva menos comisiones.</li>
                                <li>Se aplicará una modificación de la reputación del huésped.</li>
                            </ul>
                        </ol>
                        <h3>Cancelaciones por parte de anfitriones</h3>
                        <ol>
                            <li>Rango Verde: Se aplicará un reembolso total del monto de la reserva para el huésped, menos comisiones.</li>
                            <li>Rango Amarillo</li>
                            <ul>
                                <li>Se aplicará al anfitrión un cargo del 20% del monto de la reserva con un mínimo de U$D 50.</li>
                                <li>Se aplicará un reembolso total del monto de la reserva para el huésped, menos comisiones.</li>
                            </ul>
                            <li>Rango Rojo</li>
                            <ul>
                                <li>Se aplicará al anfitrión un cargo del 30% del monto de la reserva con un mínimo de U$D 50.</li>
                                <li>Se aplicará un reembolso total del monto de la reserva para el huésped, menos comisiones.</li>
                                <li>Se aplicará una modificación de la reputación del anfitrión.</li>
                            </ul>
                        </ol>
                        <p>Nota: Se podrán revisar reclamos/apelaciones de cancelaciones en circunstancias atenuantes como emergencias médicas, desastres naturales, entre otros.</p>
                    </ContentCard>
                )}
            </SmartContainer>
        </div>
    )
}

export default CancellationsPage
