package repository

import (
	"dao_vote/internal/models"
	"errors"
	"sync"
)

// Карта силы голосов для различных кошельков
var voteMap = map[string]int{
	"d016nnqrut83vd0p4afp6546rma6g5d8aqy6t7cfp": 1287398,
	"d012z4fvdrdlwp54lhqgay2j37hyk93z8ce44w5sx": 432000,
	"d01jcn7x5x4th40mjwzpcr6x6r42ws2a9kp9w22ju": 869985,
	"d014s8qq258jpcyf84ry8p6cm40mgm8v6ta49meqs": 757041,
	"d010dq9t7j93ytw7y3q9ker6460j4jx28x4vzgwgh": 486343,
	"d01hyjz895asr8mqe892nqte2dmjucq3akw8wj090": 456980,
	"d010aqhavtkn48e53dys3jcg52nphfucs7374nzxv": 342993,
	"d01p6dd9pwgxu6qx58hwy3xrctsran8u06rgn7nu4": 197101,
	"d01avgdzdya9hj8r6f37m8aslwz2rtu4lm6292ycu": 110944,
	"d01gvztujhu86pd0wyrzj8pgzm6d2n856hk5gcdjd": 611944,
	"d01xe5j36j6r5jpaxvsxl2ukz5g6vur7j7kr4jggp": 356143,
	"d017nq3vmgg4cxqq99lne6g39uya8xr75n4ny76an": 57474,
	"d01karcjx9uyer0dyla9nhe4dpljpff37ynxr8y8v": 53874,
	"d012ycamh7t20kjag30kpmefqzn5ttjzdwxq43awm": 135049,
	"d018jxc6z62n4pxsyjkxnkjvd88arhupnynepmmu0": 432000,
	"d01f99xvpyyxuna3r7z8tg4jq97nkur6rf8x0k243": 685621,
	"d019q8npegpgxw3kpqjkxhlqzncszmfxsf77nnwl2": 129000,
	"d01dk8zepstnyng2s5ua23z5lwetjrn3j65tznqt7": 782000,
	"d01a4278xvmal75y9u0sz3k9spf0vh73slk73gyf7": 432000,
	"d01xpfrc8g99j76wcgfky5vh03gr8ptknklex4std": 432000,
	"d01fpdfgg2u3pylqp02dg7tepu625k076mggx9f68": 1135500,
	"d01n0t7hl93kmnytl8yk7rhft0fhv6ragfwn3k6w7": 432000,
	"d013y0yenq7fjjdwdvg9kqsggn47mp9xh9n2ascfh": 440325,
	"d01lzp09tjydwh96euu4dv5nz2q8cucayq3l7a3we": 666179,
	"d012p9ajq6f25ndp35lfnwzdlvzyr8lp3qfer6ujm": 432000,
	"d012qlckku5v5t5kjfwg6deelyfsvxmp6rx524mca": 432000,
	"d019v3jl8wktasygh85ynjq8rzr3hfr7dvdxlcprq": 475932,
	"d01um2dvn64ref0f9nyu6g2h8xcdk4zn47td2ne5h": 58174,
	"d01ewwvevsura7d0vendwj0wm59nzzqwlhsgyjkug": 273176,
	"d01cw9mx4kpqv3m80eyrpgt65kqfn9wj94ume326w": 447605,
	"d01d4tdwssk8362hkh4a4uyn3yv0gjz7afvsgzjt9": 334000,
	"d01yxt60y8hx5ldhmplk536ujqe9p38mvhn3x2r2t": 46598,
	"d01ysxjkd8ytkem5nw4pq72cefaa7ch0dyyzfhy0p": 58079,
	"d01x09tjruh3awq59zwu4eahhzac5vlxlwht879wl": 108000,
	"d01r6fga8ahzpztxnrx69v9pxumkga8dx5vwg2xed": 235176,
	"d01w3z8zueanfyxkzwrj7yz08mylquxkgy6g7k4nx": 118000,
	"d01lyudxdpn3lkjgaa74k5h9teq5plemd04pxsrty": 108000,
	"d01ym4h7kkky0ge0s4fg43cn3lqae7mpkf2e4pcz5": 108000,
	"d015sde37tes35k3s9e9reghke2rv26twnnc5mlxp": 111777,
	"d01vftyf5vfmztn53jvt6mypkkw66e35eap0cppd0": 113307,
	"d010m3n2wyra95n63qer9yutdckdssjcta0f69zth": 166722,
	"d01kguwxafgxl80ywuk0qnk09squtzu7a42vl9p48": 169162,
	"d014rzs8y3cvjzdpr65ewaw5ly8jrt35z5f9gr35j": 53874,
	"d01ne7kvq5d9nhlv0r0jgs7v5rt6pm3k6xqs2h7vh": 2654487,
	"d01q5fgtatvvxmj7gwy7ayh0kh2rd73xkt33ma2pr": 53874,
	"d01rts9qrdt9w20wqzg3cwtc588x9gngwp6p7gw79": 53874,
	"d018pwdwguga6w4376z0mkt5lmkafy667kqptcgex": 54374,
	"d01yzmjkg9vktsyz8sd4r5mwar0w0vx54yvy0l0ms": 87197,
	"d0196ffjugts8ra7s29nufnla303lhryfshh4256w": 57174,
	"d01kscw40t0r9kdasjgwvfg9u6t34xkk2lmx8xzt0": 69000,
	"d014kqm27nfgu2qfk73ud2t4fv0zdkq5d9rsqt9zn": 72029,
	"d01ppw4jtx4wfd66fdzmauhrjsgy975upqx8dj20s": 31232,
	"d01nhz2negfv6ldpnxkacu6k6ntt4r69ukchucpus": 56927,
	"d01dvgzvth8ndu6cdayxpgwhl5kqtdk88ghwnycf0": 45000,
	"d01yhtmvy65nr92ga7lleer4lrp28ragy6xp48yyk": 30295,
	"d01ez00ry2mj3yqwgj6z6r5dk34q722rhc8zjxk8m": 45183,
	"d010uuxz0wnc9x9zhd80x9kp9urghke9n4m6espn6": 32589,
	"d01juva4qeqjyavwaf4s2vfzpg2y8vj6gl9dtne45": 1000000,
}

// totalVoices представляет общее количество голосов
var totalVoices = 12876983

// GetVoteStrength возвращает силу голоса для указанного кошелька
func GetVoteStrength(from string) (int, error) {
	strength, exists := voteMap[from]
	if !exists {
		return 0, errors.New("сторонний голос")
	}
	return strength, nil
}

// GetTotalVoices возвращает общее количество голосов
func GetTotalVoices() int {
	return totalVoices
}

// GetVoteMap возвращает карту всех голосов
func GetVoteMap() map[string]int {
	return voteMap
}

// Временное хранилище для пользовательских голосований
var (
	votes  = make(map[int]models.Vote) // Словарь для хранения голосований
	mu     sync.Mutex                  // Мьютекс для управления доступом к хранилищу
	nextID = 1                         // ID для следующего голосования
)

// SaveVote сохраняет новое пользовательское голосование
func SaveVote(vote models.Vote) (int, error) {
	mu.Lock()         // Захватываем мьютекс
	defer mu.Unlock() // Освобождаем мьютекс после выполнения функции
	vote.ID = nextID  // Присваиваем новое ID голосованию
	votes[nextID] = vote
	nextID++ // Увеличиваем ID для следующего голосования
	return vote.ID, nil
}

// GetVoteByID возвращает пользовательское голосование по его ID
func GetVoteByID(id int) (models.Vote, error) {
	mu.Lock()
	defer mu.Unlock()
	vote, exists := votes[id]
	if !exists {
		return models.Vote{}, errors.New("голосование не найдено")
	}
	return vote, nil
}

// DeleteVote удаляет пользовательское голосование по его ID
func DeleteVote(id int) error {
	mu.Lock()
	defer mu.Unlock()
	if _, exists := votes[id]; !exists {
		return errors.New("голосование не найдено")
	}
	delete(votes, id)
	return nil
}

// Временное хранилище для голосов пользователей
var (
	userVotes  = make(map[int][]models.UserVote)
	voteMu     sync.Mutex
	userVoteID = 1
)

// AddUserVote сохраняет новый голос пользователя
func AddUserVote(vote models.UserVote) (int, error) {
	voteMu.Lock()
	defer voteMu.Unlock()

	vote.VoterID = userVoteID
	userVotes[vote.VoteID] = append(userVotes[vote.VoteID], vote)
	userVoteID++
	return vote.VoterID, nil
}

// GetUserVotes возвращает все голоса пользователей для указанного голосования
func GetUserVotes(voteID int) ([]models.UserVote, error) {
	voteMu.Lock()
	defer voteMu.Unlock()

	votes, exists := userVotes[voteID]
	if !exists {
		return nil, errors.New("No user votes found for this vote")
	}
	return votes, nil
}
