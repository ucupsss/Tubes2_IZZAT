import flamme from "../assets/flamme.jpg";
import fern from "../assets/fern.jpg"
import himmel from "../assets/himmel.jpg"


const team = [
  { name: "Alya Nur Rahmah", desc: "Frontend developer", img: flamme,yapping:"Duh bingung yapping apa tapi pengen deh binge watching marathon watch list hikds when yh tubes beres udah banyak nabung manga manhwa manhua nih \n Pengen punya pabrik es krim biar bisa makan es krim tiap hari terus pengen punya bakery lucu sih baking gitu mungkin side job kalau stress sama bidang ini \n\n Ehm anu jadi kenalin aku Alya yang mendesain dan membuat tampilan website ini terus yang Himmel itu panggilannya Ucup dia yang buat DOM parser, algoritma BFS DFS dan graph untuk traversal terus yang Fern itu panggilannya Rere dia membuat web scrapping dan CSS selector"},
  { name: "Yusuf Faishal", desc: "Backend developer", img:himmel, yapping: "Empat setengah tahun difitnah-fitnah saya diam. \n Dijelek-jelekin saya juga diam. \n Dicela direndah-rendahkan saya juga diam. \n Dihujat dihujat-hujat dihina-hina saya juga diam. \n Tetapi hari ini saya sampaikan \n SAYA AKAN LAWAN!!!. "},
  { name: "Reysha Syafitri", desc: "Backend developer", img:fern, yapping: "You know you love me (Yo), I know you care (Uh-huh) \n Just shout whenever (Yo), and I'll be there (Uh-huh) \n You are my love (Yo), you are my heart (Uh-huh) \n And we will never, ever, ever be apart (Yo, uh-huh) \n Are we an item? (Yo) Girl, quit playin' (Uh-huh) \n \"We're just friends,\" (Yo) what are you sayin'? (Uh-huh) \n Said, \"There's another,\" (Yo) and looked right in my eyes (Uh-huh) \n My first love broke my heart for the first time (Yo), and I was like (Uh-huh)"},
];

export default function About() {
  return (
    <>
      <h1 className="page-title">Perkenalan</h1>
      <div className="card">
        <p>Perkenalkan kami dari kelompok IZZAT yang beranggotakan Ucup, Alya, dan Rere (gambar di bawah hanya pemanis dan kami yapping aja)</p>
      </div>
      <h1 className="page-title"></h1>
      <div className="about-grid">
        {team.map((member) => (
          <div key={member.name} className="profile">
            <img src={member.img}
              alt={member.name}
            />
            <h3>{member.name}</h3>
            <h5>{member.desc}</h5>
            <br></br>
            <p>{member.yapping}</p>
          </div>
        ))}
      </div>
    </>
  );
}