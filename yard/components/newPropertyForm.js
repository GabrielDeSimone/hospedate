import FormInputIcon from "./formInputIcon";
import FormButton from "./formButton";
import toast from 'react-hot-toast';
import utils from "../utils/utils";
import FormSelect from "./formSelect";
import PriceCalculator from "./priceCalculator";
import {useState} from "react";
const apiErrors = utils.apiErrors

function extractRoomId(url) {
    let regex = /rooms\/(\d+)/;
    let match = url.match(regex);
    return match ? match[1] : null;
}

function handleSubmit(callback) {
    async function innerHandleSubmit(event) {
        event.preventDefault()

        const roomId = extractRoomId(event.target.airbnb_room_id.value)
        if (!roomId) {
            toast.error("El link de Airbnb es inválido")
            throw new Error("Invalid airbnb link")
        }

        const endpoint = '/api/properties';
        const createResponse = await fetch(endpoint, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                maxGuests: parseInt(event.target.maxGuests.value, 10),
                city: event.target.city.value,
                price: parseInt(event.target.price.value, 10),
                airbnb_room_id: roomId,
            }),
        })
        const result = await createResponse.json()

        if (createResponse.status !== 201) {
            throw result.error
        }
        return result
    }

    return async (event) => {
        toast.loading('Creando propiedad', {
            id: 'propertyCreation'
        });
        innerHandleSubmit(event).then(() => {
            toast.dismiss('propertyCreation')
            if (callback) {
                callback()
            }
        }).catch((err) => {
            if (err === apiErrors.PropertyAlreadyTaken) {
                toast.error('La propiedad seleccionada ya ha sido creada', {
                    id: 'propertyCreation'
                })
            } else {
                toast.error('Hubo un problema al crear la propiedad', {
                    id: 'propertyCreation'
                })
            }
        })
    }
}

const NewPropertyForm = (props) => {

    const [price, setPrice] = useState(0)

    const classes = [
        props.extraClasses || '',
    ]

    return (
        <div className={classes.join(' ')}>
            <h1>Agregar nueva propiedad</h1>
            <form className="flex flex-col items-center" onSubmit={handleSubmit(props.onSubmit)}>
                <FormInputIcon icon="bed" name="airbnb_room_id" id="airbnb_room_id" placeholder="Link de Airbnb" required />
                <FormInputIcon icon="group" type="number" name="maxGuests" id="maxGuests" placeholder="Capacidad de personas" required />
                <FormSelect icon="location_on" name="city" id="city" required>
                    <option disabled="disabled" value="">Ciudad</option>
                    {utils.cities.map(city => (
                        <option key={city} value={city}>{city}</option>
                    ))}
                </FormSelect>
                <FormInputIcon icon="attach_money" type="number" name="price" id="price" value={price} onChange={(e) => {setPrice(e.target.value)}} placeholder="Precio por noche (USDT)" required />
                <ul className="list-disc px-20 text-sm text-gray-500 mt-4">
                    <li>Éste será el precio visible por los huéspedes.</li>
                    <li>Para reservas directas, coordinarás el pago con el huésped sin pagar comisiones.</li>
                    <li>Para reservas protegidas, este precio contempla la comisión del servicio (saber más).</li>
                </ul>
                <PriceCalculator price={price} />
                <div className="mt-8 w-[80%] sm:w-auto flex justify-center">
                    <FormButton type="submit" text="Crear propiedad" />
                </div>
            </form>
        </div>
    )
}

export default NewPropertyForm