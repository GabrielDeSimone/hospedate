import utils from "../utils";

const PropertyInfoBar = (props) => {

    let items = [
        {
            icon: "location_on",
            key: "city",
            text: () => props.property.city,
        },
        {
            icon: "group",
            key: "maxGuests",
            text: () => `Admite hasta ${props.property.maxGuests} huéspedes`,
        },
        {
            icon: "attach_money",
            key: "price",
            text: () => `${props.property.price} USDT por noche`
        },

        // remove check-in and check-out items if the owner is seeing the page
        ...(props.isOwner ? [] : [
            {
                icon: "date_range",
                key: "checkinDate",
                text: () => `Check-in el ${utils.humanizeDate(props.checkinDate)}`
            },
            {
                icon: "date_range",
                key: "checkoutDate",
                text: () => `Check-out el ${utils.humanizeDate(props.checkoutDate)}`
            }
        ])
    ]

    utils.amenitiesList.map(amenity => {
        if (props.property[amenity.key] !== null) {
            items.push({
                icon: amenity.icon,
                key: amenity.key,
                text: amenity.isBoolean ? (() => amenity.strText) : (() => {
                    return amenity.strFull(props.property[amenity.key])
                }),
            })
        }
    })


    return (
        /*<div className="grid grid-cols-1 tablet:grid-cols-2 desktop:grid-cols-3">*/
        <ul className="grid grid-cols-1 tablet:grid-cols-2 desktop:grid-cols-3">
            {items.map(item => (
                <li key={item.key} className="flex flex-row px-3 my-2">
                    <span className="
                            material-symbols-outlined
                            text-2xl
                            bg-transparent
                            text-gray-700
                            pt-1
                        ">{item.icon}</span>
                    <div className="mt-1.5 ml-1.5">{item.text()}</div>
                </li>
            ))}
        </ul>
    )
}

export default PropertyInfoBar