import NavBar from '../../components/navBar'
import {useContext, useEffect, useState} from 'react'
import { UserContext } from '../../components/userProvider'
import Router, {useRouter} from "next/router";
import Link from "next/link";
import utils from "../../components/utils";
import toast from "react-hot-toast";
import {TableH, TableHTd, TableHTh, TableHTr} from "../../components/myTable";
import FormButton from "../../components/formButton";


function VisitorEphemeral({ reservation }) {
    return (
        <div>
            <h1 className="text-2xl md:mb-7">Reserva # {reservation.id} (Pendiente de pago)</h1>
            <div className="bg-[#fcf8d5] p-10 rounded-xl">
                <h2 className="font-bold">⚠ Tu pago está pendiente - Quedan aproximadamente {utils.humanizeRemainingTime(utils.remainingEphemeralSeconds(reservation.createdAt))}</h2>
                <div className="bg-[#fffef3] p-4 rounded-xl">
                    <p>Debes hacer uno o varios pagos por el monto de $ {utils.humanizeCents(reservation.totalBilledCents)} USDT a la direción que se muestra a continuación para que tu reserva sea comunicada al anfitrión.</p>
                </div>
                <h3 className="text-lg mt-3">Dirección de pago</h3>
                <TableH>
                    <TableHTr>
                        <TableHTh>Moneda</TableHTh>
                        <TableHTd>USDT</TableHTd>
                    </TableHTr>
                    <TableHTr>
                        <TableHTh>Red</TableHTh>
                        <TableHTd>TRON - TRC-20</TableHTd>
                    </TableHTr>
                    <TableHTr>
                        <TableHTh>Dirección</TableHTh>
                        <TableHTd>{reservation.walletAddress} <Link href={`https://tronscan.org/#/address/${reservation.walletAddress}`} target="_blank" className="ml-2">Haz click aquí para ver el balance actual</Link></TableHTd>
                    </TableHTr>
                </TableH>
                <div className="bg-[#fffef3] p-4 rounded-xl mt-4">
                    <ul className="list-disc ml-7">
                        <li>Antes de realizar un envío de dinero, comprueba que la dirección es correcta. Hospedate no puede ayudarte si envías tu dinero por error a otro destinatario.</li>
                        <li>Si ya realizaste tu pago, vuelve en unos minutos.</li>
                        <li>Si no cuentas con fondos para pagar en la moneda y red especificadas, puedes <Link href="#" onClick={() => cancelReservation(reservation)}>cancelar tu reserva</Link>.</li>
                    </ul>
                </div>
            </div>
            {reservationInfo(reservation)}
        </div>
    )
}

function VisitorPending({ reservation }) {
    return (
        <div>
            <h1 className="text-2xl md:mb-7">Reserva # {reservation.id}</h1>
            <div className="bg-[#e3ecff] p-10 rounded-xl">
                <h2 className="font-bold">¡Tu reserva está pendiente!</h2>
                <div className="bg-[#FBFCFF] p-4 rounded-xl">
                    <p>Nos hemos comunicado con el anfitrión. Te avisaremos cuando la reserva haya sido confirmada.</p>
                </div>
            </div>
            {reservationInfo(reservation)}
        </div>
    )
}

function VisitorConfirmed({ reservation }) {
    return (
        <div>
            <h1 className="text-2xl md:mb-7">Reserva # {reservation.id} (Confirmada)</h1>
            <h3>Información del anfitrión</h3>
            {showOwnerInfo(reservation.owner)}
            <h3>Información de la reserva</h3>
            {reservationInfo(reservation)}
            <div className="flex justify-center">
                <FormButton text="Cancelar reserva" onClick={() => {cancelReservation(reservation)}}/>
            </div>
        </div>
    )
}

function VisitorInProgress({ reservation }) {
    return (
        <div>
            <h1 className="text-2xl md:mb-7">Reserva # {reservation.id} (Confirmada)</h1>
            <h3>Información del anfitrión</h3>
            {showOwnerInfo(reservation.owner)}
            <h3>Información de la reserva</h3>
            {reservationInfo(reservation)}
            <p className="text-gray-500 mt-3 ml-8 text-sm">La reserva se encuentra en progreso.</p>
        </div>
    )
}

function VisitorCanceled({ reservation }) {
    let canceledMsg
    if (reservation.canceledBy === utils.CANCELED_BY_OWNER) {
        canceledMsg = "(Cancelada por el anfitrión)"
    } else {
        canceledMsg = "(Cancelada)"
    }
    return (
        <div>
            <h1 className="text-2xl md:mb-7">Reserva # {reservation.id} {canceledMsg}</h1>
            {reservationInfo(reservation)}
        </div>
    )
}

function VisitorDiscarded({ reservation }) {
    return (
        <div>
            <h1 className="text-2xl md:mb-7">Reserva # {reservation.id} (Descartada)</h1>
            {reservationInfo(reservation)}
        </div>
    )
}

function showOwnerInfo(userOwner) {
    return (
        <TableH>
            <TableHTr>
                <TableHTh>Nombre</TableHTh>
                <TableHTd>{userOwner.name}</TableHTd>
            </TableHTr>
            <TableHTr>
                <TableHTh>Correo electrónico</TableHTh>
                <TableHTd>{userOwner.email}</TableHTd>
            </TableHTr>
        </TableH>
    )
}

function showVisitorInfo(userVisitor) {
    return (
        <TableH>
            <TableHTr>
                <TableHTh>Nombre</TableHTh>
                <TableHTd>{userVisitor.name}</TableHTd>
            </TableHTr>
            <TableHTr>
                <TableHTh>Correo electrónico</TableHTh>
                <TableHTd>{userVisitor.email}</TableHTd>
            </TableHTr>
            <TableHTr>
                <TableHTh>Número de teléfono</TableHTh>
                <TableHTd>{userVisitor.phoneNumber}</TableHTd>
            </TableHTr>
        </TableH>
    )
}

async function confirmReservation(reservation) {
    toast.loading("Confirmando reserva", {
        id: "confirmingReservation"
    })
    const confirmReq = await fetch(
        `/api/reservations/${reservation.id}/confirm`,
        {
            method: "POST",
        }
    )
    if (confirmReq.status === 200) {
        toast.success("La reserva fue confirmada", {
            id: "confirmingReservation"
        })
        setTimeout(() => {
            window.location.reload()
        }, 1000)
    } else {
        toast.error("Hubo un problema al confirmar la reserva", {
            id: "confirmingReservation"
        })
    }
}

async function cancelReservation(reservation) {
    toast.loading("Cancelando reserva", {
        id: "cancellingReservation"
    })
    const cancelReq = await fetch(
        `/api/reservations/${reservation.id}/cancel`,
        {
            method: "POST",
        }
    )
    if (cancelReq.status === 200) {
        toast.success("La reserva fue cancelada", {
            id: "cancellingReservation"
        })
        setTimeout(() => {
            window.location.reload()
        }, 1000)
    } else {
        toast.error("Ocurrió un problema al cancelar la reserva", {
            id: "cancellingReservation"
        })
    }
}

function OwnerPending({ reservation, userVisitor }) {
    return (
        <div>
            <h1 className="text-2xl md:mb-7">Solicitud de Reserva # {reservation.id} (Pendiente de aceptación)</h1>
            <h2>Información del huésped</h2>
            {showVisitorInfo(userVisitor)}
            <p className="text-gray-500 mt-3 ml-8 text-sm">Te recomendamos contactarte con el huésped antes de confirmar la reserva.</p>
            <h2 className="mt-8">Información de la reserva</h2>
            {reservationInfo(reservation)}
            <p className="text-gray-500 mt-3 ml-8 text-sm">Transcurridas 5 horas de la creación de la reserva, será cancelada automáticamente si no has respondido.</p>
            <div className="flex justify-center">
                <FormButton text="Cancelar reserva" onClick={() => {cancelReservation(reservation)}} />
                <FormButton text="Aceptar reserva" extraClasses="ml-4" onClick={() => {confirmReservation(reservation)}} />
            </div>
        </div>
    )
}

function OwnerConfirmed({ reservation }) {
    return (
        <div>
            <h1 className="text-2xl md:mb-7">Reserva # {reservation.id}</h1>
            <h2 className="mt-8">Información de la reserva</h2>
            {reservationInfo(reservation)}
            <div className="flex justify-center">
                <FormButton text="Cancelar reserva" onClick={() => {cancelReservation(reservation)}}/>
            </div>
        </div>
    )
}

function OwnerInProgress({ reservation }) {
    return (
        <div>
            <h1 className="text-2xl md:mb-7">Reserva # {reservation.id}</h1>
            <h2 className="mt-8">Información de la reserva</h2>
            {reservationInfo(reservation)}
            <div className="flex justify-center">
                <FormButton text="Cancelar reserva" onClick={() => {cancelReservation(reservation)}}/>
            </div>
        </div>
    )
}

function OwnerCompleted({ reservation }) {
    return (
        <div>
            <h1 className="text-2xl md:mb-7">Reserva # {reservation.id}</h1>
            <h2 className="mt-8">Información de la reserva</h2>
            {reservationInfo(reservation)}
            <p className="text-gray-500 mt-3 ml-8 text-sm">La reserva ha finalizado.</p>
        </div>
    )
}

function OwnerCanceled({ reservation, userVisitor }) {
    let canceledMsg
    if (reservation.canceledBy === utils.CANCELED_BY_VISITOR) {
        canceledMsg = "(Cancelada por el huésped)"
    } else {
        canceledMsg = "(Cancelada)"
    }
    return (
        <div>
            <h1 className="text-2xl md:mb-7">Reserva # {reservation.id} {canceledMsg}</h1>
            <h2>Información del huésped</h2>
            {showVisitorInfo(userVisitor)}
            <h2 className="mt-8">Información de la reserva</h2>
            {reservationInfo(reservation)}
        </div>
    )
}

function OwnerDiscarded({ reservation, userVisitor }) {
    return (
        <div>
            <h1 className="text-2xl md:mb-7">Reserva # {reservation.id} (descartada)</h1>
            <h2>Información del huésped</h2>
            {showVisitorInfo(userVisitor)}
            <h2 className="mt-8">Información de la reserva</h2>
            {reservationInfo(reservation)}
        </div>
    )
}

function reservationInfo(reservation) {
    return (
        <TableH>
            <TableHTr>
                <TableHTh>Estado</TableHTh>
                <TableHTd>{utils.humanizeReservationStatus(reservation.status)}</TableHTd>
            </TableHTr>
            <TableHTr>
                <TableHTh>Propiedad</TableHTh>
                <TableHTd><Link href={`/properties/${reservation.propertyId}`} className="text-hosp-light-blue hover:underline cursor-pointer" target="_blank">Ver propiedad</Link></TableHTd>
            </TableHTr>
            <TableHTr>
                <TableHTh>Fecha de check-in</TableHTh>
                <TableHTd>{utils.humanizeFullDate(reservation.checkinDate)} {remainingDaysLabel(reservation.checkinDate)}</TableHTd>
            </TableHTr>
            <TableHTr>
                <TableHTh>Fecha de check-out</TableHTh>
                <TableHTd>{utils.humanizeFullDate(reservation.checkoutDate)}</TableHTd>
            </TableHTr>
            <TableHTr>
                <TableHTh>Forma de pago</TableHTh>
                <TableHTd>{utils.humanizeReservationType(reservation.reservationType)}</TableHTd>
            </TableHTr>
            <TableHTr>
                <TableHTh>Precio total</TableHTh>
                <TableHTd>USDT $ {utils.humanizeCents(reservation.totalBilledCents)}</TableHTd>
            </TableHTr>
        </TableH>
    )
}





function remainingDaysLabel(targetDateStr) {
    const daysDifference = utils.getDaysDifference(new Date(), targetDateStr)
    if (daysDifference >= 0) {
        return `(Quedan ${daysDifference} días)`
    } else {
        return ""
    }
}

async function fetchReservation(reservationId) {
    const getReservation = await fetch(
        `/api/reservations/${reservationId}`
    )
    const response = await getReservation.json()

    if (response.data) {
        return response.data
    } else if (getReservation.status === 404 || getReservation.status === 403) {
        Router.push('/404')
    } else {
        return {}
    }
}

function whichComponent(reservation, user) {
    const isOwner = reservation.ownerId === user.id

    if (isOwner && reservation.status === utils.RESERVATION_STATUS_EPHEMERAL) {
        return (() => (<div>owner ephemeral</div>))()
    } else if (isOwner && reservation.status === utils.RESERVATION_STATUS_PENDING) {
        return <OwnerPending reservation={reservation} userVisitor={reservation.user} />
    } else if (isOwner && reservation.status === utils.RESERVATION_STATUS_CONFIRMED) {
        return <OwnerConfirmed reservation={reservation} />
    } else if (isOwner && reservation.status === utils.RESERVATION_STATUS_IN_PROGRESS) {
        return <OwnerInProgress reservation={reservation} />
    } else if (isOwner && reservation.status === utils.RESERVATION_STATUS_CANCELED) {
        return <OwnerCanceled reservation={reservation} userVisitor={reservation.user} />
    } else if (isOwner && reservation.status === utils.RESERVATION_STATUS_DISCARDED) {
        return <OwnerDiscarded reservation={reservation} />
    } else if (isOwner && reservation.status === utils.RESERVATION_STATUS_COMPLETED) {
        return <OwnerCompleted reservation={reservation} />
    } else if (!isOwner && reservation.status === utils.RESERVATION_STATUS_EPHEMERAL) {
        return <VisitorEphemeral reservation={reservation} />
    } else if (!isOwner && reservation.status === utils.RESERVATION_STATUS_PENDING) {
        return <VisitorPending reservation={reservation} />
    } else if (!isOwner && reservation.status === utils.RESERVATION_STATUS_CONFIRMED) {
        return <VisitorConfirmed reservation={reservation} />
    } else if (!isOwner && reservation.status === utils.RESERVATION_STATUS_IN_PROGRESS) {
        return <VisitorInProgress reservation={reservation} />
    } else if (!isOwner && reservation.status === utils.RESERVATION_STATUS_CANCELED) {
        return <VisitorCanceled reservation={reservation} />
    } else if (!isOwner && reservation.status === utils.RESERVATION_STATUS_DISCARDED) {
        return <VisitorDiscarded reservation={reservation} />
    } else if (!isOwner && reservation.status === utils.RESERVATION_STATUS_COMPLETED) {
        return (() => (<div>visitor completed</div>))()
    } else {
        return (() => (<div>unknown</div>))()
    }
}


function Reservation() {
    const [user, setUser] = useContext(UserContext);
    const [reservation, setReservation] = useState({})
    const router = useRouter();

    useEffect(() => {
        (async () => {
            toast.loading("Cargando reserva", {
                id: "loadingReservation"
            })
            if (router.query.id) {
                const reserv = await fetchReservation(router.query.id)
                if (reserv && reserv.id) {
                    toast.dismiss("loadingReservation")
                    setReservation(reserv)

                    if (! utils.reservationStatusIsFinal(reserv.status)) {
                        setTimeout(() => {
                            window.location.reload()
                        }, 120000)
                    }
                }
            }
        })()
    }, [router.isReady])

    return (
        <div>
            <NavBar />
            <div className="w-full m-auto bg-white px-5 md:px-10 py-5 mt-4">
                {reservation && reservation.id && user && user.id && whichComponent(reservation, user)}
            </div>
        </div>
    );
}

export default Reservation