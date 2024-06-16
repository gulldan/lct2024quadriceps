import {
  IconBrandGithub,
  IconBrandLinkedin,
  IconBrandTwitter,
} from "@tabler/icons-react";
import Image from "next/image";
import React from "react";

const TEAM = [
  {
    name: "Алексей Дубынин",
    img: "/team/lehad.jpg",
    description: `Алексей Дубынин is currently a standing cat.`,
    github: "https://github.com/sevenineone",
    linkedin: "https://youtu.be/dQw4w9WgXcQ",
    twitter: "https://youtu.be/dQw4w9WgXcQ",
  },
  {
    name: "Алексей Галиев",
    img: "/team/lehag.jpeg",
    description: `Алексей Галиев is a rust cat.`,
    github: "https://github.com/TiregeRRR",
    linkedin: "https://youtu.be/dQw4w9WgXcQ",
    twitter: "https://youtu.be/dQw4w9WgXcQ",
  },

  {
    name: "Рустам Сахабутдинов",
    img: "/team/rustam.jpg",
    description: `Рустам Сахабутдинов is a cat with legs.`,
    github: "https://github.com/RustamOper05",
    linkedin: "https://youtu.be/dQw4w9WgXcQ",
    twitter: "https://youtu.be/dQw4w9WgXcQ",
  },
  {
    name: "Максим Клещенок",
    img: "/team/max.jpg",
    description: `Максим Клещенок - Bullet cat.`,
    github: "https://github.com/gulldan",
    linkedin: "https://youtu.be/dQw4w9WgXcQ",
  },
  {
    name: "Эльвин Шахбазов",
    img: "/team/elvin.jpg",
    description: `I am jerry`,
    github: "https://github.com/semwai",
    linkedin: "https://youtu.be/dQw4w9WgXcQ",
  },
];

export const MENTORS = [
  {
    name: "Юлия Гумилёва",
    img: "/team/yulenka.jpeg",
    description: "помогла нам попасть в правильные ритмы",
    linkedin: "facebook.com/julia.gumilova",
  },
  {
    name: "Валерий Балашов",
    img: "/team/valeryb.jpg",
    description: "BO$$. СПАСИБО ЗА ДЕТСТВО",
    linkedin: "https://t.me/Leg1of",
  },
];

export default function Gallery() {
  return (
    <>
      <section className="space-y-16 md:space-y-32 my-14">
        {TEAM.map((element, index) => (
          <div
            className={`flex flex-col justify-evenly ${
              index % 2 === 0 ? "md:flex-row" : "md:flex-row-reverse"
            }`}
            key={index}
          >
            <Image
              className="lg:w-[20%] md:w-[30%] w-[60%] aspect-w-4 aspect-h-3 rounded-2xl self-center"
              src={element.img}
              width={500}
              height={300}
              alt={`Picture of ${element.name}`}
            />
            <div className="md:w-[40%] mt-2 md:mt-0">
              <h2 className="text-4xl mb-4 font-medium">{element.name}</h2>
              <p className="mb-4 text-gray-600 md:leading-9 md:text-lg">
                {element.description}
              </p>
              <div className="flex self-end gap-4">
                <a
                  target="_blank"
                  href={element.github}
                  alt={`github profile link of ${element.name}`}
                  className="hover:text-primary-900 transition-all duration-300"
                >
                  <IconBrandGithub />
                </a>
                <a
                  target="_blank"
                  href={element.linkedin}
                  className="hover:text-primary-900 transition-all duration-300"
                  alt={`linkedin profile link of ${element.name}`}
                >
                  <IconBrandLinkedin />
                </a>
              </div>
            </div>
          </div>
        ))}
      </section>
      {/* Mentor Section */}
      <section className="space-y-16 md:space-y-32 my-14">
        <h1 className="text-5xl font-bold text-center mt-12 pt-12">
          Our Mentor
        </h1>
        {MENTORS.map((mentor, index) => (
          <div
            className={`flex flex-col justify-evenly ${
              index % 2 === 0 ? "md:flex-row" : "md:flex-row-reverse"
            }`}
            key={index}
          >
            <Image
              className="lg:w-[20%] md:w-[30%] w-[60%] aspect-w-4 aspect-h-3 rounded-2xl self-center"
              src={mentor.img}
              width={500}
              height={300}
              alt={`Picture of ${mentor.name}`}
            />
            <div className="md:w-[40%] mt-2 md:mt-0">
              <h2 className="text-4xl mb-4 font-medium">{mentor.name}</h2>
              <p className="mb-4 text-gray-600 md:leading-9 md:text-lg">
                {mentor.description}
              </p>
              <div className="flex self-end gap-4">
                {mentor.linkedin && (
                  <a
                    target="_blank"
                    href={mentor.linkedin}
                    className="hover:text-primary-900 transition-all duration-300"
                    alt={`linkedin link of ${mentor.name}`}
                  >
                    <IconBrandLinkedin />
                  </a>
                )}
                {mentor.twitter && (
                  <a
                    target="_blank"
                    href={mentor.twitter}
                    className="hover:text-primary-900 transition-all duration-300"
                  >
                    <IconBrandTwitter />
                  </a>
                )}
              </div>
            </div>
          </div>
        ))}
      </section>
    </>
  );
}
