import { Swiper, SwiperSlide } from 'swiper/react';
import { Navigation, Pagination, Keyboard } from 'swiper';
import 'swiper/css';
import 'swiper/css/navigation';
import 'swiper/css/pagination';

const Slider = (props) => {

    const classes = [
        props.extraClasses || '',
    ]

    return (
        <Swiper
            modules={[Keyboard, Navigation, Pagination]}
            navigation
            keyboard={{
                enabled: true,
            }}
            pagination={{ clickable: true }}
            spaceBetween={10}
            slidesPerView={1}
            onSwiper={(swiper) => console.log(swiper)}
        >
            {props.images && props.images.map(image => (
                <SwiperSlide key={image} className="h-full">
                    <img src={image} className="object-contain h-[500px] m-auto" />
                </SwiperSlide>
            ))}
        </Swiper>
    )
}

export default Slider