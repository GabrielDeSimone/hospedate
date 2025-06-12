import PropertyCard from "./propertyCard"

const PropertyGrid = (props) => {

    const properties = props.properties || [];

    const originalGridClasses = [
        "grid",
        "grid-cols-1",
        "md:grid-cols-2",
        "min-[690px]:grid-cols-2",
        "min-[950px]:grid-cols-3",
        "min-[1240px]:grid-cols-4",
        "min-[1640px]:grid-cols-5",
        "min-[1880px]:grid-cols-6",
        "gap-4",
    ]

    const classes = [
        "mt-6",
        props.extraClasses || '',
    ]

    const gridClasses = props.customGridClasses ? props.customGridClasses : originalGridClasses

    return (
        <ul className={classes.join(' ') + ' ' + gridClasses.join(' ')}>
                {
                properties.map(property =>
                    <li className="grid-child" key={property.id}>
                        <PropertyCard
                            id={property.id}
                            title={property.title}
                            price={property.price}
                            maxGuests={property.maxGuests}
                            images={property.images}
                            status={property.status}

                            stayParams={props.stayParams}
                            />
                    </li>)
                }
        </ul>
    )
}

export default PropertyGrid