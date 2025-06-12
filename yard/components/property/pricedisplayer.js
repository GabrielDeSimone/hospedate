import MakeOrderButton from "./makeOrderButton";
import React, {useEffect, useState} from "react";
import Modal from "../modal";
import utils from "../utils";
import FormButton from "../formButton";
import toast from 'react-hot-toast';
import Router from "next/router";

const RESERVATION_HOSPEDATE = "in_platform"
const RESERVATION_DIRECT_OWNER = "owner_directly"

function handleReservation(reservationType, stayParams, propertyId) {
    async function innerHandle(event) {
        event.preventDefault()
        toast.loading('Reservando', {
            id: 'makingReservation'
        })
        const makeReservation = await fetch(
            "/api/reservations",
            {
                method: "POST",
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    checkinDate: stayParams.checkinDate,
                    checkoutDate: stayParams.checkoutDate,
                    guests: stayParams.guests,
                    reservationType: reservationType,
                    propertyId: propertyId
                }),
            }
        )
        const response = await makeReservation.json()
        if (response.data && response.data.id) {
            toast.success('La reserva se hizo exitosamente', {
                id: 'makingReservation'
            })
            setTimeout(() => {
                Router.push(`/reservations/${response.data.id}`)
            }, 1000)
        } else {
            toast.error('Hubo un problema al hacer la reserva', {
                id: 'makingReservation'
            })
        }
    }

    return innerHandle
}

function generateConfirmationContent(reservationType, stayParams, totalPrice, propertyId) {
    const paymentMethod = () => {
        if (reservationType === RESERVATION_HOSPEDATE) {
            return (<li>Método de pago: Transferencia de USDT via red TRON (TRC-20)</li>)
        } else {
            return (<li>Método de pago: A convenir con el anfitrión</li>)
        }
    }
    return (
        <>
            <h2>Confirmación de reserva</h2>
            <form onSubmit={handleReservation(reservationType, stayParams, propertyId)}>
                <ul>
                    <li>Fecha de check-in: {utils.humanizeFullDate(stayParams.checkinDate)}</li>
                    <li>Fecha de check-out: {utils.humanizeFullDate(stayParams.checkoutDate)}</li>
                    <li>Cantidad de noches: {stayParams.stayDays}</li>
                    <li>Cantidad de huéspedes: {stayParams.guests}</li>
                    {paymentMethod()}
                    <li>Precio final: USDT $ {totalPrice}</li>
                </ul>
                <FormButton type="submit" text="Reservar" />
            </form>
        </>
    )
}


const PriceDisplayer = (props) => {

    const [modalOpen, setModalOpen] = useState(false);
    const [confirmationContent, setConfirmationContent] = useState((<div></div>));
    const [totalPrice, setTotalPrice] = useState('-');

    const handleModalClose = () => {
        setModalOpen(false);
    }

    const handleModalOpen = (reservationType) => {
        setConfirmationContent(generateConfirmationContent(reservationType, props.stayParams, totalPrice, props.propertyId))
        setModalOpen(true)
    }

    useEffect(() => {
        const totalPriceComp = (() => {
            if (props.stayParams && props.price) {
                return props.stayParams.stayDays * props.price
            } else {
                return '-'
            }
        })()
        setTotalPrice(totalPriceComp)
    }, [props.price, props.stayParams])

    return (
        <div>
            <div className="flex"> {/*per night price*/}
                <span className="mt-1.5">Precio por noche:</span>
                <div className="flex ml-5">
                    <span className="text-2xl">$ {props.price || '-'}</span>
                    <span className="ml-2 text-hospedate-green">USDT</span>
                </div>
            </div>
            <div className="mt-4"> {/*total price*/}
                <span>Precio total por {props.stayParams.stayDays || '-'} noches:</span>
                <div className="flex justify-center mt-5">
                    <span className="text-5xl font-light">$ {totalPrice}</span>
                    <span className="ml-2 text-2xl text-hospedate-green">USDT</span>
                </div>
            </div>
            <div className="rounded-xl mt-5 border-gray-300 border p-5"> {/*hospedate order*/}
                <ul className="max-w-xs text-gray-700">
                    <li><span className="text-hospedate-green">✓</span> Pagá de forma online</li>
                    <li className="mt-2"><span className="text-hospedate-green">✓</span> Reservá con seguridad: si no podés hacer check-in, te devolvemos el dinero.</li>
                </ul>
                <MakeOrderButton extraClasses="mt-4" bgColor="bg-hospedate-green" text="Reserva Protegida" onSubmit={() => handleModalOpen(RESERVATION_HOSPEDATE)} />
            </div>
            <div className="rounded-xl mt-5 border-gray-300 border p-5"> {/*direct owner order*/}
                <ul className="max-w-xs text-gray-700">
                    <li><span className="text-hospedate-green">✓</span> Ganá flexibilidad al acordar la forma de pago más conveniente con el anfitrión.</li>
                    <li className="mt-2"><span className="text-hospedate-green">⚠</span> Hospedate no puede protegerte si arreglas de forma directa. ¡Ten cuidado!</li>
                </ul>
                <MakeOrderButton extraClasses="mt-4" bgColor="bg-hosp-pale-blue" text="Reserva Directa" onSubmit={() => handleModalOpen(RESERVATION_DIRECT_OWNER)} />
            </div>
            <Modal show={modalOpen} onClose={handleModalClose} position="top">
                {confirmationContent}
            </Modal>
        </div>
    )
}

export default PriceDisplayer