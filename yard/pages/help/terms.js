import NavBar from '../../components/navBar';
import {useContext, useState} from "react";
import {UserContext, isGuest, isLoggedIn} from "../../components/userProvider";
import Router from "next/router";
import ContentCard from "../../components/contentCard";
import Link from "next/link";
import SmartContainer from "../../components/smartContainer";

const TermsPage = () => {

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
                        <h1>Términos y Condiciones</h1>
                        <h2>1. Introducción</h2>
                            <p>Al hacer uso de la plataforma "Hospedate" (en adelante, “la plataforma”), el usuario reconoce y acepta cumplir con los Términos y Condiciones detallados a continuación.</p>
                        <h2>2. Sobre los servicios ofrecidos</h2>
                            <p>La plataforma ofrece a los usuarios un servicio donde pueden realizar reservas de alquileres temporarios. Los usuarios tienen la posibilidad de hacer reservas de tipo “directas” o “protegidas”. Con las reservas directas ambas partes se pueden contactar directamente y concretar el pago entre ellas. Las reservas protegidas se efectúan a través del sistema de pago con criptomonedas de la plataforma.</p>
                            <p>La plataforma se compromete a brindar el mejor servicio posible a sus usuarios, buscando siempre la satisfacción y seguridad en todas las transacciones y comunicaciones realizadas a través de ella.</p>
                            <p>La plataforma se encuentra actualmente en periodo de pruebas, por lo que es posible que los usuarios encuentren errores en su funcionamiento.</p>
                        <h2>3. Sobre las acciones del usuario</h2>
                            <p>La plataforma se reserva el derecho de suspender o eliminar cuentas de usuarios que, a juicio discrecional de Hospedate, no actúen de buena fe, intenten engañar a la plataforma o a otros usuarios, o violen cualquiera de los términos aquí expuestos.</p>
                            <p>La plataforma asume que todas las comunicaciones vía email, mensajes de texto o mensajes de WhatsApp entre la plataforma y la dirección de email y número de teléfonos proporcionados por el usuario son de autoría del usuario y legítimas. La plataforma no se responsabiliza si el usuario pierde el control de su dirección de email o número de teléfono, ni por la pérdida de fondos o daños que puedan ocurrir en consecuencia.</p>
                        <h2>4. Sobre las reservas no confirmadas</h2>
                            <p>Los anfitriones se reservan el derecho de rechazar reservas o no confirmarlas.</p>
                            <p>Las reservas pendientes de confirmación que no sean respondidas por los anfitriones serán rechazadas automáticamente luego del plazo de 5 horas.</p>
                        <h2>5. Sobre las reservas canceladas</h2>
                            <p>Los reembolsos, retenciones, cargos, compensaciones o modificaciones en la reputación de los usuarios debido a cancelaciones de reservas están definidas en nuestra Política de Cancelaciones.</p>
                            <p>Los reembolsos a huéspedes se realizarán mediante una transacción en criptomonedas por la cantidad correspondiente según se establezca en la Política de Cancelaciones hacia una dirección proporcionada por el usuario por WhatsApp luego de confirmación con un representante de la plataforma.</p>
                        <h2>Sobre eventos no previstos</h2>
                            <p>La plataforma no se hace responsable por daños a la propiedad de parte de los huéspedes o terceros.</p>
                            <p>La plataforma no se hace responsable por engaños, estafas o fraudes ocurridas por reservas de tipo “directas”.</p>
                            <p>En caso de engaños o estafas en reservas de tipo “protegidas” la plataforma evaluará cada caso puntualmente y determinará si aplicar uno o varios reembolsos, retenciones, cargos, compensaciones o modificaciones en la reputación de los usuarios involucrados.</p>
                        <h2>7. Sobre las acreditaciones</h2>
                            <p>Las acreditaciones a las cuentas de los anfitriones tomarán lugar luego de transcurridas 36 horas de la fecha de check-out.</p>
                        <h2>8. Sobre los retiros</h2>
                            <p>La plataforma se compromete a realizar y completar todas las solicitudes de retiro de criptomonedas que los usuarios efectúen a través de la plataforma. Una vez procesada la solicitud, las cantidades nominales solicitadas serán enviadas a la dirección proporcionada por el usuario, exceptuando comisiones por uso del sistema de pagos o la red.</p>
                        <h2>9. Sobre la naturaleza de las criptomonedas y los virus informáticos</h2>
                            <p>La plataforma no se hace responsable de virus, programas informáticos dañinos o dispositivos defectuosos que el usuario pueda utilizar y que causen daños a su propiedad, a terceros, o la pérdida parcial o total de fondos en criptomonedas que el usuario haya retirado y tenga bajo su control.</p>
                            <p>La plataforma no se hace responsable por las fluctuaciones en el valor de las criptomonedas utilizadas en la plataforma. Sin embargo, garantiza que todos los retiros de criptomonedas en cantidades nominales sean realizados y completados según lo solicitado por el usuario.</p>
                            <p>La plataforma no realiza ninguna inversión ni movimientos innecesarios con los saldos en criptomonedas que los usuarios mantienen en la plataforma. Los fondos permanecen estáticos hasta que el usuario decida hacer uso de ellos o solicite su retiro.</p>
                        <h2>10. Sobre normativas e impuestos</h2>
                            <p>La plataforma recomienda a sus usuarios conocer las jurisdicciones que les sean aplicables y cumplir con las normativas o reglamentaciones municipales o gubernamentales que correspondan. A su vez la plataforma no se hace responsable por incumplimientos a las mismas.</p>
                        <h2>11. Sobre los Términos y Condiciones</h2>
                            <p>La plataforma se reserva el derecho de modificar los Términos y Condiciones de uso en cualquier momento y sin previo aviso. Es responsabilidad del usuario revisar periódicamente estos términos.</p>
                            <p>Al hacer uso de la plataforma, el usuario reconoce y acepta cumplir con los Términos y Condiciones.</p>
                        <h2>12. Jurisdicción y Legislación Aplicable</h2>
                            <p>Cualquier disputa o conflicto relacionado con la interpretación o ejecución de estos Términos y Condiciones se resolverá de conformidad con las leyes y regulaciones vigentes en el estado de Delaware, Estados Unidos, donde está registrada la sede principal de Hospedate. Las partes acuerdan someterse a la jurisdicción exclusiva de los tribunales de Delaware para la resolución de cualquier disputa que pueda surgir en relación con estos términos y condiciones.</p>
                        <h2>13. Contacto</h2>
                            <p>Para cualquier consulta o aclaración sobre estos Términos y Condiciones, los usuarios pueden contactar al equipo de soporte de Hospedate a través de los medios proporcionados en la plataforma.</p>
                    </ContentCard>
                )}
            </SmartContainer>
        </div>
    )
}

export default TermsPage
