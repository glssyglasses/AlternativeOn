//Version: 0.0.5 (Beta 5)
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gonutz/w32/v2"
	"golang.org/x/sys/windows/registry"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/ncruces/zenity"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type AutoGenerated []struct {
	ID       string      `json:"id"`
	Nome     string      `json:"nome"`
	Cor      string      `json:"cor"`
	Ordem    int         `json:"ordem"`
	Ativo    bool        `json:"ativo"`
	Solucoes []Solutions `json:"solucoes"`
}
type Solutions struct {
	ID               string `json:"id"`
	Nome             string `json:"nome"`
	Descricao        string `json:"descricao"`
	Arquivo          string `json:"arquivo"`
	Link             string `json:"link"`
	Ativo            bool   `json:"ativo"`
	TipoRenderizacao string `json:"tipoRenderizacao"`
	Slug             string `json:"slug"`
	Ordem            int    `json:"ordem"`
	DataCadastro     string `json:"dataCadastro"`
	Novo             bool   `json:"novo"`
}

type Error struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

var mainwin *ui.Window

func init() {
	//Check on registry if command prompt should be visible
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Princess Mortix\Alternative On`, registry.QUERY_VALUE)
	if err != nil {
		fmt.Println("[w] Caminho da configuração do registro não encontrado, crindo...")
		k, _, err = registry.CreateKey(registry.CURRENT_USER, `Software\Princess Mortix\Alternative On`, uint32(registry.CURRENT_USER))
		if err != nil {
			fmt.Println("[w] Erro ao criar chave no registro:", err, "\nContinuando mesmo assim...")
		}
		fmt.Println("[i] Configuração criada com sucesso.")
		k.Close()
		return
	}
	defer k.Close()
	val, _, err := k.GetIntegerValue("HideCmd")
	if err != nil {
		fmt.Println("[i] O terminal não está configurado para ser ocultado, configurando...")
		setKey, err := registry.OpenKey(registry.CURRENT_USER, `Software\Princess Mortix\Alternative On`, registry.SET_VALUE)
		if err != nil {
			fmt.Println("[w] Erro ao criar configuração para o terminal:", err, "\nContinuando mesmo assim...")
		}
		err = setKey.SetDWordValue("HideCmd", 1)
		if err != nil {
			fmt.Println("[w] Erro ao configurar o terminal:", err, "\nContinuando mesmo assim...")
		}
		return
	}
	if val == 1 {
		//1 is true
		console := w32.GetConsoleWindow()
		if console != 0 {
			_, consoleProcID := w32.GetWindowThreadProcessId(console)
			if w32.GetCurrentProcessId() == consoleProcID {
				w32.ShowWindowAsync(console, w32.SW_HIDE)
			}
		}
	}
	k.Close()
	zenity.Warning("Este é a beta de um cliente alternativo ao Positivo On, bugs podem ocorrer", zenity.Title("Aviso"), zenity.WarningIcon)

	fmt.Println("[i] Initialization complete")
}

func loginPage() ui.Control {
	//Create a login window on a vertical box
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	vbox.Append(ui.NewLabel("Para começarmos, faça seu login. Use o usuário e senha do Positivo On"), false)
	vbox.Append(ui.NewVerticalSeparator(), false)

	group := ui.NewGroup("Login")
	group.SetMargined(true)
	vbox.Append(group, true)

	aboutform := ui.NewForm()
	aboutform.SetPadded(true)
	group.SetChild(aboutform)

	aboutform.Append("Usuário", ui.NewEntry(), false)
	aboutform.Append("Senha", ui.NewPasswordEntry(), false)
	loginButton := ui.NewButton("Entrar")
	loginButton.OnClicked(func(*ui.Button) {
		//Show a information dialog
		ui.MsgBox(mainwin, "Login", "Esse botão não funciona ainda, mas você pode fechar essa janela para fazer login.")

	})
	aboutform.Append("", loginButton, false)

	return vbox
}

func aboutPage() ui.Control {
	//Create an about window on a vertical box
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	vbox.Append(ui.NewLabel("Aqui você pode saber sobre o projeto Alternative On."), false)
	vbox.Append(ui.NewVerticalSeparator(), false)

	group := ui.NewGroup("Sobre o Alternative On")
	group.SetMargined(true)
	vbox.Append(group, true)

	aboutform := ui.NewForm()
	aboutform.SetPadded(true)
	group.SetChild(aboutform)

	//add text on a new multi-line entry read only
	abouttext := ui.NewMultilineEntry()
	abouttext.SetReadOnly(true)
	abouttext.SetText("Este é um cliente alternativo ao Positivo On, ele foi desenvolvido para facilitar o uso da plataforma, e ainda está na fase beta.\n\nO projeto é mantido pelo grupo de desenvolvimento da Princess Mortix, e é um projeto open source, você pode acessar o projeto no github clicando no botão abaixo.\nCaso você queira contribuir com o projeto, você pode acessar o github você pode enviar um issue, ou se preferir, você pode enviar um pull request diretamente no github.\n\nObrigado por usar o Alternative On!\n\n**Politica de Privacidade:**\nNós não coletamos nenhum dado de você, mas talvez a plataforma Positivo On, que é um serviço de terceiros, coleta dados de usuários.")
	aboutform.Append("", abouttext, true)
	//add a button to open the github page
	ghlink := ui.NewButton("Acessar o projeto no github")
	ghlink.OnClicked(func(*ui.Button) {
		//open the github page
		err := w32.ShellExecute(0, "open", "https://github.com/PrincessMortix/AlternativeOn", "", "", w32.SW_SHOW)
		if err != nil {
			fmt.Println("[E]", err)
		}
	})
	terms := ui.NewButton("Termos do Positivo On")
	terms.OnClicked(func(*ui.Button) {
		//open the github page
		err := w32.ShellExecute(0, "open", "https://positivoon.com.br/#/termos-de-uso", "", "", w32.SW_SHOW)
		if err != nil {
			fmt.Println("[E]", err)
		}
	})
	privacy := ui.NewButton("Política de Privacidade do Positivo On")
	privacy.OnClicked(func(*ui.Button) {
		//open the github page
		err := w32.ShellExecute(0, "open", "https://positivoon.com.br/#/politica-de-privacidade", "", "", w32.SW_SHOW)
		if err != nil {
			fmt.Println("[E]", err)
		}
	})
	aboutform.Append("", terms, false)
	aboutform.Append("", privacy, false)
	aboutform.Append("", ghlink, false)

	return vbox
}

func settingsPage() ui.Control {
	//Create a settings window on a vertical box
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	vbox.Append(ui.NewLabel("Aqui você pode alterar o funcionamento do Alternative On"), false)
	vbox.Append(ui.NewVerticalSeparator(), false)

	group := ui.NewGroup("Configurações disponíveis")
	group.SetMargined(true)
	vbox.Append(group, true)

	//add text on a new multi-line entry read only
	settingsForm := ui.NewForm()
	settingsForm.SetPadded(true)
	group.SetChild(settingsForm)

	settingsText := ui.NewMultilineEntry()
	settingsText.SetReadOnly(true)
	settingsText.SetText("Ajuda das configurações:\n\n- Mostrar terminal de depuração (Essa configuração mostra o terminal, mostrando sobre possiveis erros da aplicação)")
	settingsForm.Append("", settingsText, true)

	settingsShowDebug := ui.NewCheckbox("Mostrar terminal de depuração")
	settingsShowDebug.OnToggled(func(*ui.Checkbox) {
		if settingsShowDebug.Checked() {
			fmt.Println("[i] O terminal foi configurado para ficar ativo.")
			//Configure the terminal to show debug messages
			setTerminal, err := registry.OpenKey(registry.CURRENT_USER, `Software\Princess Mortix\Alternative On`, registry.SET_VALUE)
			if err != nil {
				fmt.Println("[E] Erro ao criar configuração para o terminal:", err)
				ui.MsgBoxError(mainwin, "Erro", "Erro ao criar configuração para o terminal: "+err.Error()+"\nTente abrir a aplicação novamente como administrador.")
				os.Exit(1)
			}
			err = setTerminal.SetDWordValue("HideCmd", 0)
			if err != nil {
				fmt.Println("[E] Erro ao configurar o terminal:", err)
				ui.MsgBoxError(mainwin, "Erro", "Erro ao configurar o terminal: "+err.Error()+"\nTente abrir a aplicação novamente como administrador.")
				os.Exit(1)
			}
		} else {
			fmt.Println("[i] O terminal foi configurado para ficar inativo.")
			//Configure the terminal to hide debug messages
			setTerminal, err := registry.OpenKey(registry.CURRENT_USER, `Software\Princess Mortix\Alternative On`, registry.SET_VALUE)
			if err != nil {
				fmt.Println("[w] Erro ao criar configuração para o terminal:", err)
				ui.MsgBoxError(mainwin, "Erro", "Erro ao criar configuração para o terminal: "+err.Error()+"\nTente abrir a aplicação novamente como administrador.")
				os.Exit(1)
			}
			err = setTerminal.SetDWordValue("HideCmd", 1)
			if err != nil {
				fmt.Println("[w] Erro ao configurar o terminal:", err)
				ui.MsgBoxError(mainwin, "Erro", "Erro ao configurar o terminal: "+err.Error()+"\nTente abrir a aplicação novamente como administrador.")
				os.Exit(1)
			}
		}
	})
	settingsForm.Append("", settingsShowDebug, false)

	return vbox
}
func setupUI() {
	mainwin = ui.NewWindow("Alternative On", 640, 480, true)
	mainwin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})

	tab := ui.NewTab()
	mainwin.SetChild(tab)
	mainwin.SetMargined(true)

	tab.Append("Inicio", loginPage())
	tab.SetMargined(0, true)
	tab.Append("Sobre", aboutPage())
	tab.SetMargined(1, true)
	tab.Append("Configurações", settingsPage())

	mainwin.Show()
}
func main() {
	ui.Main(setupUI)
	user, err := zenity.Entry("Qual o seu usuário do Positivo On?", zenity.Title("Usuário"), zenity.InfoIcon)
	if err != nil {
		fmt.Println("[i] {username}", err)
		os.Exit(0)
	}
	if user == "" {
		zenity.Error("Usuário não informado", zenity.Title("Erro"), zenity.ErrorIcon)
		fmt.Println("[w] Usuário não informado")
		main()
	}
	_, pass, err := zenity.Password(zenity.Title("Qual a sua senha do Positivo On?"))
	if err != nil {
		fmt.Println("[i] {password}", err)
		os.Exit(0)
	}
	if pass == "" {
		zenity.Error("Senha não informada", zenity.Title("Erro"), zenity.ErrorIcon)
		fmt.Println("[w] Senha não informada")
		main()
	}
	fmt.Println("[d] Sending initial request")

	//Post request and unmarshal response.
	token, err := postReq("https://sso.specomunica.com.br/connect/token", "POST", "username="+user+"&password="+pass+"&grant_type=password&client_id=hubpsd&client_secret=DA5730D8-90FF-4A41-BFED-147B8E0E2A08&scope=openid%20offline_access%20integration_info")
	if err != nil {
		fmt.Println("[E] Login falhou", err)
		os.Exit(1)
	}
	fmt.Println("[d] Login sucess!\n[d - 166] Token:", token)
	solutions, err := authReq("https://apihub.positivoon.com.br/api/Categoria/Solucoes/Perfil/ALUNO?NivelEnsino=EM&IdEscola=149259e8-6864-4f41-a3c8-6c624184bc56", "GET", token)
	if err != nil {
		fmt.Println(err)
	}
	//Unmarshal response
	fmt.Println("[d] Solutions loaded sucessfully!!\n[d - 122]:", solutions)
	zenity.Info("As soluções foram carregadas com sucesso!", zenity.Title("Soluções Carregadas!"), zenity.InfoIcon)
	//Unmarshal json response
	var autoGenerated AutoGenerated
	err = json.Unmarshal([]byte(solutions), &autoGenerated)
	if err != nil {
		fmt.Println(err)
	}

	//Add solutions to list
	var list []string
	for _, solution := range autoGenerated {
		list = append(list, solution.Nome)
	}
	//Print solutions list
	fmt.Println("[d] Solutions list:", list)
	//Get links from solutions
	var links []string
	for _, solution := range autoGenerated {
		for _, sol := range solution.Solucoes {
			links = append(links, sol.Link)
		}
	}
	//Print links
	fmt.Println("[d] Links:", links)

	//Ask for solution
	ssolution, err := zenity.List("Escolha uma solução", list, zenity.Title("Solução"))
	fmt.Println("[d] Solution:", ssolution)
	if err != nil {
		fmt.Println(err)
	}
	//switch solution
	switch ssolution {
	case "Avaliação":
		//Open Avaliações in browser
		//First, get link from Avaliações from json
		var links []string
		for _, solution := range autoGenerated {
			for _, sol := range solution.Solucoes {
				links = append(links, sol.Link)
			}
		}
		//Print links
		fmt.Println("[d] {switch} Links:", links)
		//Print solution
		fmt.Println("[d] {switch} Solution:", ssolution)
		//change {token} to token
		newlink := strings.Replace(links[0], "{token}", token, -1)
		fmt.Println("[d] {switch} New link:", newlink)
		//Open link in browser
		err := w32.ShellExecute(0, "open", newlink, "", "", w32.SW_SHOW)
		if err != nil {
			fmt.Println("[E]", err)
		}

	default:
		fmt.Println("[d] {switch} Solution:", ssolution)
		zenity.Warning("Solução não implementada ainda, tente em uma próxima build :)", zenity.Title("Solução não implementada"), zenity.WarningIcon)
	}

}

func postReq(url string, method string, payload string) (string, error) {

	client := &http.Client{}
	fmt.Println("[d | postReq] Payload:", payload)
	req, err := http.NewRequest(method, url, strings.NewReader(payload))

	if err != nil {
		fmt.Println("[e | postReq]", err)
		return "Failed to create request:", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("[e | postReq]", err)
		return "Failed to send request: ", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("[e | postReq]", err)
		return "Failed to read response: ", err
	}

	//fmt.Println("[d] Response:", string(body))
	//Check if response is valid
	if res.StatusCode != 200 {
		//Unmarshal error response
		var errResp Error
		json.Unmarshal(body, &errResp)
		fmt.Println("[e | postReq] Error:", errResp.Error)
		fmt.Println("[e | postReq] Error description:", errResp.ErrorDescription)
		zenity.Error(errResp.ErrorDescription, zenity.Title("Erro"), zenity.ErrorIcon)
		return errResp.Error + ": " + errResp.ErrorDescription, err
	}
	//Unmarshal response
	var token Token
	json.Unmarshal(body, &token)
	//convert int to string
	expiresIn := fmt.Sprintf("%d", token.ExpiresIn)
	fmt.Println("[d | postReq] Token: "+token.AccessToken, "\n[d | postReq] Refresh Token: "+token.RefreshToken, "\n[d | postReq] Expires In: "+expiresIn, "\n[d | postReq] Token Type: "+token.TokenType)
	zenity.Info("Login realizado com sucesso!\n\nIremos agora tentar carregar as soluções...", zenity.Title("Login realizado com sucesso!"), zenity.InfoIcon)
	return token.AccessToken, nil
}

func authReq(url string, method string, token string) (string, error) {

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("[e | authReq]", err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("[e | authReq]", err)
		return "", err
	}
	return string(body), nil
}
